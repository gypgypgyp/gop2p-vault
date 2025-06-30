package store

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"bytes"
)

// HashKey generates a SHA-1 hash for a file stream
func HashKey(r io.Reader) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}


// HashKeyBytes returns SHA-1 hash of raw byte data
func HashKeyBytes(data []byte) (string, error) {
	return HashKey(bytes.NewReader(data))
}

// Optional helper for compatibility
func BytesReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}