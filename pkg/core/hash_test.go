package core

import (
	"sync"
	"testing"
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
			result := ComputeHash(tt.name, tt.seed, tt.seq)
			resultHex := make([]byte, len(result)*2)
			for i, b := range result {
				resultHex[i*2] = "0123456789abcdef"[b>>4]
				resultHex[i*2+1] = "0123456789abcdef"[b&0x0f]
			}
			if string(resultHex) != tt.expectedHex {
				t.Errorf("expected %s, got %s", tt.expectedHex, string(resultHex))
			}
		})
	}
}

func TestHasherPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = ComputeHash("test", []byte("seed"), uint64(i))
		}(i)
	}
	wg.Wait()
}
