package rng

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateLocal(t *testing.T) {
	tests := []struct {
		name string
		size int
		err  bool
	}{
		{"Generate 16 bytes", 16, false},
		{"Generate 32 bytes", 32, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call GenerateLocal with the specified size.
			seed, err := GenerateLocal(tt.size)
			require.NoError(t, err, "GenerateLocal should not return an error")
			require.NotNil(t, seed, "GenerateLocal should return a non-nil seed")
			require.Len(t, seed, tt.size, "Generated seed should have the correct length")
			t.Logf("Generated seed: %x", seed)
		})
	}
}
