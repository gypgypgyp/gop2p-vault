package store

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

// HashKey generates a SHA-1 hash for a file stream
func HashKey(r io.Reader) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
