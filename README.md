# seedselection-go

**WARNING: This project is experimental and work-in-progress, DO NOT USE IN PRODUCTION.**

![Build & Test (Core)](https://github.com/spacecomputerio/seedselection-go/actions/workflows/build_test.yml/badge.svg?branch=main)

## Overview

This package provides a simple and effective ways to implement deterministic selection in distributed systems. 

The focus is on random seeds to provide fair selection and good distribution. The selection process is deterministic, meaning that the same input will always produce the same output. This is useful for ensuring that all peers in a distributed (and trustless) network agree on the selected items, without introducing network calls or the need for a central authority.

A max-heap is used to efficiently select the top `n` items based on a distance function (xor distance) applied to the each of items against the hash of seed+seq+name.

### Security Considerations

- The selection depends on the entropy of the seed, so it is crucial to use a secure and unpredictable seed to ensure fairness and security. We encourage users to use SpaceComputer's [orbitport](https://docs.spacecomputer.io/orbitport) as a source of randomness.
- The hash function is pluggable, it should be cryptographically secure to prevent manipulation of the selection process.

## Usage

```go
package main

import (
    "crypto/sha256"
    "fmt"
    "log"

	seedselection "github.com/spacecomputerio/seedselection-go/pkg/core"
    // seedselection_rng "github.com/spacecomputerio/seedselection-go/pkg/rng"
)

func main() {
	// rand, _ := seedselection_rng.GenerateLocal(32)
	rand, _ := hex.DecodeString("dd59b31bb729d255426208cc99df1d7a3f233e6b86013d47df77aab88a204baa")

	peerIDS := []string{
		"698750a09b934337746f0973448167f364cae132e2f8b327ae4913e5b5445029",
		"3b213ced003e89b35a26c22cbd011c9bfab29578415b2069f7fc8b01998b903d",
		"e42bbf8533f4f0b1d44e7fc1c9ac54a6ac368642dd1b8a10a1775255eed0c31a",
		"a7a0243e04fd71dc10068134a7dc0ab6de6e3cb76439400d17e6d531a5e596b1",
		"b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o",
		"c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o3p4",
	}
	round := uint64(1)
	n := 2
	// Create a new peer selection instance
	selected, err := seedselection.XorDistanceSelection(sha256.New(), "testgroup", rand, round, n, peerIDS)
	if err != nil {
		log.Fatalf("Failed to create peer selection: %v", err)
	}

	fmt.Printf("%+v\n", selected)
    // Output: [698750a09b934337746f0973448167f364cae132e2f8b327ae4913e5b5445029 3b213ced003e89b35a26c22cbd011c9bfab29578415b2069f7fc8b01998b903d]
}
```

## License

This project is licensed under the terms of the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions to this project! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

