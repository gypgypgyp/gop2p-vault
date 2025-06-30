package store

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"bytes"
	// "fmt"
	"path/filepath"
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

// HashPath converts a SHA-1 hash string into a sharded file path like ./data/66/c5/<hash>
func HashPath(baseDir, hash string) string {
	sub1 := hash[0:2]
	sub2 := hash[2:4]
	return filepath.Join(baseDir, sub1, sub2, hash)
}