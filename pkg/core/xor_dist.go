package core

import (
	"container/heap"
	"hash"
	"math/big"
)

// XorDistanceSelection takes a name, (random) seed, sequence number, number of ids to select (n), and a list of ids (32 byte hashes) to select from.
// The hash(name, seed, seq, n) will be used for ordering the items according to the xor distance, where the closest n ids will be selected.
// For optimization, a max-heap is used to maintain the n ids with the smallest xor distances from the hash value.
// If n is greater than the number of ids, all ids will be elected.
func XorDistanceSelection(hasher hash.Hash, name string, seed []byte, seq uint64, n int, ids []string) ([]string, error) {
	p := len(ids)
	if p == 0 {
		return nil, nil
	}
	if n >= p {
		return ids, nil
	}
	hash := ComputeHash(hasher, name, seed, seq)
	hashValue := new(big.Int).SetBytes(hash)

	maxHeap := newMaxDistanceHeap()
	// creating reusable big.Ints to hold the id and xor distance values
	itemValue := new(big.Int)
	distanceVal := new(big.Int)
	for _, item := range ids {
		itemValue.SetBytes([]byte(item))
		distanceVal.Xor(hashValue, itemValue)
		maxHeap.PushUpToN(heapEntry{
			id:       item,
			distance: distanceVal.Uint64(),
		}, n)
	}

	return maxHeap.List(), nil
}

// heapEntry is an internal struct that holds an id and its xor distance from a hash value.
type heapEntry struct {
	id       string
	distance uint64
}

// maxDistanceHeap is a max-heap of heapEntry.
// It is used to maintain the n ids with the smallest xor distances from a given hash value,
// enabling efficient selection of the closest ids based on xor distance.
type maxDistanceHeap []heapEntry

func newMaxDistanceHeap() *maxDistanceHeap {
	h := &maxDistanceHeap{}
	heap.Init(h)
	return h
}

func (h maxDistanceHeap) Len() int {
	return len(h)
}

// For a max-heap, we want elements with LARGER 'distance' to be "less" (to be pushed to the top)
func (h maxDistanceHeap) Less(i, j int) bool {
	return h[i].distance > h[j].distance
}
func (h maxDistanceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *maxDistanceHeap) Push(x any) {
	*h = append(*h, x.(heapEntry))
}
func (h *maxDistanceHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// PushUpToN adds a new heapEntry to the maxDistanceHeap, ensuring that the heap does not exceed size n.
// If the heap is already at size n, it will only add the new entry if its distance is smaller than the largest distance in the heap.
// This maintains the invariant that the heap contains the n ids with the smallest distances.
// This is useful for selecting the closest ids based on xor distance.
func (h *maxDistanceHeap) PushUpToN(entry heapEntry, n int) {
	if h.Len() < n {
		heap.Push(h, entry)
	} else {
		last := (*h)[0] // Get the largest element in the max-heaps
		// If the current distance is smaller than the largest distance in the heap (root of max-heap)
		if entry.distance < last.distance {
			heap.Pop(h)         // Remove the current largest
			heap.Push(h, entry) // Add the new smaller one
		}
	}
}

// List will clean the max-heap and return the ids in the order they were selected
// i.e. from closest to furthest
func (h *maxDistanceHeap) List() []string {
	ids := make([]string, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		entry := heap.Pop(h).(heapEntry)
		ids[i] = entry.id
	}
	return ids
}
