package lib

import (
	"crypto/rand"
	"fmt"
)

func Random(size uint) ([]byte, error) {
	out := make([]byte, size)

	if _, err := rand.Read(out); err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}

	return out, nil
}
