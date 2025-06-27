package core

import (
	"hash"
)

// Options holds the config parameters for the selection process.
type Options struct {
	// weights is an optional parameter that allows for weighted selection.
	weights []int64
	// hasher is the hash function to use for computing the xor distance.
	hasher hash.Hash
}

type Opt func(opts *Options)

// WithWeights allows for weighted selection of ids.
func WithWeights(weights []int64) Opt {
	return func(opts *Options) {
		opts.weights = weights
	}
}

// WithHasher allows for specifying a custom hash function to use for computing the xor distance.
func WithHasher(hasher hash.Hash) Opt {
	return func(opts *Options) {
		opts.hasher = hasher
	}
}

func NewOptions(opts ...Opt) *Options {
	o := new(Options)
	for _, opt := range opts {
		opt(o)
	}
	return o
}
