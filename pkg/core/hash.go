package core

import (
	"crypto/sha256"
	"hash"
	"strconv"
	"sync"
)

// ComputeHash computes a hash based on the name, seed and sequence number / round / epoch.
func ComputeHash(name string, seed []byte, seq uint64) []byte {
	hasher := getHasher()
	defer putHasher(hasher)

	hasher.Write([]byte(name))
	hasher.Write(seed)
	hasher.Write([]byte(strconv.Itoa(int(seq))))

	return hasher.Sum(nil)
}

// pool for sha256 hasher to avoid frequent allocations
// and deallocations, which can be expensive in high-load scenarios.
// It is used to get a hasher for computing the hash of the peer selection function.
var hasherPool = sync.Pool{
	New: func() interface{} {
		return sha256.New()
	},
}

// getHasher retrieves a sha256 hasher from the pool.
func getHasher() hash.Hash {
	return hasherPool.Get().(hash.Hash)
}

// putHasher returns a sha256 hasher to the pool after resetting it.
func putHasher(h hash.Hash) {
	h.Reset()
	hasherPool.Put(h)
}
