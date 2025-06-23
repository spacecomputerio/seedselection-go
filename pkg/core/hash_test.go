package core

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeHash(t *testing.T) {
	tests := []struct {
		testName    string
		name        string
		seed        []byte
		seq         uint64
		expectedHex string
	}{
		{
			testName:    "Happy path",
			name:        "test",
			seed:        []byte("seed"),
			seq:         1,
			expectedHex: "468460ee3c32ca9574f91f213853d0b0aece116aa74b71ab66bb7a9c558b2b7c",
		},
		{
			testName:    "No name",
			name:        "",
			seed:        []byte("seed"),
			seq:         1,
			expectedHex: "df9ecf4c79e5ad77701cfc88c196632b353149d85810a381f469f8fc05dc1b92",
		},
		{
			testName:    "Empty seed",
			name:        "test",
			seed:        []byte{},
			seq:         1,
			expectedHex: "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result := ComputeHash(sha256.New(), tt.name, tt.seed, tt.seq)
			resultHex := hex.EncodeToString(result)
			require.Equal(t, tt.expectedHex, string(resultHex), "Hash mismatch")
		})
	}
}

func TestHasherPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = ComputeHash(nil, "test", []byte("seed"), uint64(i))
		}(i)
	}
	wg.Wait()
}
