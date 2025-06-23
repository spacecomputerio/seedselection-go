package core

import (
	"crypto/sha256"
	"hash"
	"strconv"
	"sync"
)

// ComputeHash computes a hash based on the name, seed and sequence number / round / epoch.
// It accepts a hash.Hash interface, which allows for flexibility in the hashing algorithm used.
// If no hasher is provided, it uses a sha256 hasher from a sync.Pool to avoid frequent allocations.
func ComputeHash(hasher hash.Hash, name string, seed []byte, seq uint64) []byte {
	if hasher == nil {
		hasher = getHasher()
		defer putHasher(hasher)
	}

	hasher.Write([]byte(name))
	hasher.Write(seed)
	hasher.Write([]byte(strconv.Itoa(int(seq))))

	return hasher.Sum(nil)
}

// pool for sha256 hasher to avoid frequent allocations
// and deallocations, which can be expensive in high-load scenarios.
// It is used to get a hasher for computing the hash of the seed selection function.
var sha256HasherPool = sync.Pool{
	New: func() interface{} {
		return sha256.New()
	},
}

// getHasher retrieves a sha256 hasher from the pool.
func getHasher() hash.Hash {
	return sha256HasherPool.Get().(hash.Hash)
}

// putHasher returns a sha256 hasher to the pool after resetting it.
func putHasher(h hash.Hash) {
	h.Reset()
	sha256HasherPool.Put(h)
}
