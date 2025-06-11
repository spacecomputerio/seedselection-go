package core

import (
	"container/heap"
	"math/big"
)

// XorDistanceSelection takes a name, (random) seed, sequence number, number of peer to select (n), and a list of peers (32 byte hashes/ids) to select from.
// The hash(name, seed, seq, n) will be used for ordering the peerset according to the xor distance, where the closest n peers will be selected.
// For optimization, a max-heap is used to maintain the n peers with the smallest xor distances from the hash value.
// If n is greater than the number of peers, all peers will be elected.
func XorDistanceSelection(name string, seed []byte, seq uint64, n int, peerset []string) ([]string, error) {
	p := len(peerset)
	if p == 0 {
		return nil, nil
	}
	if n >= p {
		return peerset, nil
	}
	hash := ComputeHash(name, seed, seq)
	hashValue := new(big.Int).SetBytes(hash)

	maxHeap := newMaxPeerDistanceHeap()
	// creating reusable big.Ints to hold the peer id and xor distance values
	peerValue := new(big.Int)
	distanceVal := new(big.Int)
	for _, peer := range peerset {
		peerValue.SetBytes([]byte(peer))
		distanceVal.Xor(hashValue, peerValue)
		maxHeap.PushUpToN(heapEntry{
			peer:     peer,
			distance: distanceVal.Uint64(),
		}, n)
	}

	return maxHeap.Peerset(), nil
}

// heapEntry is an internal struct that holds a peer and its xor distance from a hash value.
type heapEntry struct {
	peer     string
	distance uint64
}

// maxPeerDistanceHeap is a max-heap of peerEntry.
// It is used to maintain the n peers with the smallest xor distances from a given hash value,
// enabling efficient selection of the closest peers based on xor distance.
type maxPeerDistanceHeap []heapEntry

func newMaxPeerDistanceHeap() *maxPeerDistanceHeap {
	h := &maxPeerDistanceHeap{}
	heap.Init(h)
	return h
}

func (h maxPeerDistanceHeap) Len() int {
	return len(h)
}

// For a max-heap, we want elements with LARGER 'distance' to be "less" (to be pushed to the top)
func (h maxPeerDistanceHeap) Less(i, j int) bool {
	return h[i].distance > h[j].distance
}
func (h maxPeerDistanceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *maxPeerDistanceHeap) Push(x any) {
	*h = append(*h, x.(heapEntry))
}
func (h *maxPeerDistanceHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// PushUpToN adds a new peerEntry to the maxPeerDistanceHeap, ensuring that the heap does not exceed size n.
// If the heap is already at size n, it will only add the new entry if its distance is smaller than the largest distance in the heap.
// This maintains the invariant that the heap contains the n peers with the smallest distances.
// This is useful for selecting the closest peers based on xor distance.
func (h *maxPeerDistanceHeap) PushUpToN(entry heapEntry, n int) {
	if h.Len() < n {
		heap.Push(h, entry)
	} else {
		last := (*h)[h.Len()-1] // Get the largest element in the max-heaps
		// If the current distance is smaller than the largest distance in the heap (root of max-heap)
		if entry.distance < last.distance {
			heap.Pop(h)         // Remove the current largest
			heap.Push(h, entry) // Add the new smaller one
		}
	}
}

// Peerset will clean the max-heap and return the peers in the order they were selected
// i.e. from closest to furthest
func (h *maxPeerDistanceHeap) Peerset() []string {
	peerset := make([]string, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		entry := heap.Pop(h).(heapEntry)
		peerset[i] = entry.peer
	}
	return peerset
}
