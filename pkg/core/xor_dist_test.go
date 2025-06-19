package core

import (
	"encoding/hex"
	"testing"

	"github.com/spacecomputerio/seedselection-go/pkg/rng"
	"github.com/stretchr/testify/require"
)

func TestXorDistanceSelection(t *testing.T) {
	peerIDS := []string{
		"peer1",
		"peer2",
		"peer3",
		"peer4",
		"peer5",
		"peer6",
		"peer7",
		"peer8",
		"peer9",
		"peer10",
	}

	tests := []struct {
		name        string
		rngName     string
		seed        []byte
		seq         uint64
		n           int
		peerset     []string
		expected    []string
		expectedErr bool
	}{
		{
			name: "seq 1 with 5 peers",
			seed: []byte("test-seed"),
			seq:  1,
			n:    3,
			peerset: func() []string {
				peerset := make([]string, 5)
				copy(peerset, peerIDS[0:5])
				return peerset
			}(),
			expected:    []string{peerIDS[4], peerIDS[3], peerIDS[0]},
			expectedErr: false,
		},
		{
			name: "seq 1 with 5 peers with different order",
			seed: []byte("test-seed"),
			seq:  1,
			n:    3,
			peerset: func() []string {
				peerset := make([]string, 5)
				copy(peerset, peerIDS[0:5])
				peerset[0], peerset[1] = peerset[1], peerset[0] // swap first two for different order
				peerset[2], peerset[3] = peerset[3], peerset[2] // swap next two for different order
				return peerset
			}(),
			expected:    []string{peerIDS[4], peerIDS[3], peerIDS[0]},
			expectedErr: false,
		},
		{
			name: "seq 2 with 5 peers",
			seed: []byte("test-seed"),
			seq:  2,
			n:    3,
			peerset: func() []string {
				peerset := make([]string, 5)
				copy(peerset, peerIDS[0:5])
				return peerset
			}(),
			expected:    []string{peerIDS[3], peerIDS[4], peerIDS[1]},
			expectedErr: false,
		},
		{
			name:    "Seq 1 with rng name",
			rngName: "dummy-rng-name-x",
			seed:    []byte("test-seed"),
			seq:     1,
			n:       3,
			peerset: func() []string {
				peerset := make([]string, 5)
				copy(peerset, peerIDS[0:5])
				return peerset
			}(),
			expected:    []string{peerIDS[3], peerIDS[4], peerIDS[0]},
			expectedErr: false,
		},
		{
			name: "N greater than peerset size",
			seed: []byte("test-seed"),
			seq:  1,
			n:    10,
			peerset: func() []string {
				peerset := make([]string, 3)
				copy(peerset, peerIDS[0:3])
				return peerset
			}(),
			expected:    []string{peerIDS[0], peerIDS[1], peerIDS[2]},
			expectedErr: false,
		},
		{
			name:        "No nodes",
			seed:        []byte("test-seed"),
			seq:         1,
			n:           0,
			peerset:     []string{},
			expected:    nil,
			expectedErr: false,
		},
		{
			name:    "with 32-byte peer IDs",
			rngName: "testgroup-1",
			seed:    []byte("test-seed"),
			seq:     10,
			n:       2,
			peerset: []string{
				"698750a09b934337746f0973448167f364cae132e2f8b327ae4913e5b5445029",
				"3b213ced003e89b35a26c22cbd011c9bfab29578415b2069f7fc8b01998b903d",
				"e42bbf8533f4f0b1d44e7fc1c9ac54a6ac368642dd1b8a10a1775255eed0c31a",
				"a7a0243e04fd71dc10068134a7dc0ab6de6e3cb76439400d17e6d531a5e596b1",
				"b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o",
				"c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o3p4",
			},
			expected: []string{
				"e42bbf8533f4f0b1d44e7fc1c9ac54a6ac368642dd1b8a10a1775255eed0c31a",
				"a7a0243e04fd71dc10068134a7dc0ab6de6e3cb76439400d17e6d531a5e596b1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := tt.rngName
			if name == "" {
				name = "test"
			}
			actual, err := XorDistanceSelection(name, tt.seed, tt.seq, tt.n, tt.peerset)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, actual)
			}
		})
	}
}

// generatePeerset generates a peerset of 'p' random 32-byte hash strings.
func generatePeerset(p int) []string {
	peerset := make([]string, p)
	for i := 0; i < p; i++ {
		p, err := rng.GenerateLocal(32)
		if err != nil {
			panic("failed to generate random peer ID: " + err.Error())
		}
		peerset[i] = hex.EncodeToString(p)

	}
	return peerset
}

// benchmarkXorDistanceSelection is a helper to reduce boilerplate in benchmark functions.
func benchmarkXorDistanceSelection(b *testing.B, p, n int) {
	peerset := generatePeerset(p)

	for i := 0; b.Loop(); i++ {
		_, err := XorDistanceSelection("benchmark_service", []byte("benchmark_election_seed"), uint64(i), n, peerset)
		if err != nil {
			b.Fatalf("XorDistanceSelection failed: %v", err)
		}
	}
}

// --- Benchmarks for varying N (selected leaders) and P (total peers) ---

func BenchmarkXorDistanceSelection_N10_P100(b *testing.B) {
	benchmarkXorDistanceSelection(b, 100, 10)
}

func BenchmarkXorDistanceSelection_N10_P1000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 1000, 10)
}

func BenchmarkXorDistanceSelection_N10_P10000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 10000, 10)
}

func BenchmarkXorDistanceSelection_N10_P100000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 100000, 10)
}

func BenchmarkXorDistanceSelection_N10_P500000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 500000, 10)
}

func BenchmarkXorDistanceSelection_N10_P1000000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 1000000, 10)
}

func BenchmarkXorDistanceSelection_N100_P1000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 1000, 100)
}

func BenchmarkXorDistanceSelection_N100_P10000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 10000, 100)
}

func BenchmarkXorDistanceSelection_N100_P100000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 100000, 100)
}

func BenchmarkXorDistanceSelection_N100_P500000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 500000, 100)
}

func BenchmarkXorDistanceSelection_N100_P1000000(b *testing.B) {
	benchmarkXorDistanceSelection(b, 1000000, 100)
}
