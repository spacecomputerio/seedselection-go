package rng

import (
	gocryptorand "crypto/rand"
)

// GenerateLocal generates a new random seed of the specified size.
// It uses the system's secure random number generator to ensure the seed is cryptographically secure.
// The size parameter specifies the number of bytes to generate.
func GenerateLocal(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := gocryptorand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
