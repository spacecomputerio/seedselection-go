package main

import (
	"encoding/hex"
	"fmt"
	"log"

	seedselection "github.com/spacecomputerio/seedselection-go/pkg/core"
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
	selected, err := seedselection.XorDistanceSelection("testgroup", rand, round, n, peerIDS)
	if err != nil {
		log.Fatalf("Failed to create peer selection: %v", err)
	}

	fmt.Printf("%+v\n", selected)
	// Output: [698750a09b934337746f0973448167f364cae132e2f8b327ae4913e5b5445029 3b213ced003e89b35a26c22cbd011c9bfab29578415b2069f7fc8b01998b903d]

	// Verify the selected items
	expected := []string{
		"698750a09b934337746f0973448167f364cae132e2f8b327ae4913e5b5445029",
		"3b213ced003e89b35a26c22cbd011c9bfab29578415b2069f7fc8b01998b903d",
	}
	if len(selected) != len(expected) {
		log.Fatalf("Expected %d builders, got %d", len(expected), len(selected))
	}
	for i, builder := range selected {
		if builder != expected[i] {
			log.Fatalf("Expected builder %d to be %s, got %s", i, expected[i], builder)
		}
	}
}
