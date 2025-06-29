package store

import (
	"io"
	"os"
	"path/filepath"
)

type Store struct {
	BaseDir string // Root directory for all files
}

// New creates a new Store instance
func New(baseDir string) *Store {
	return &Store{BaseDir: baseDir}
}

// Has checks whether the file with given key exists
func (s *Store) Has(key string) bool {
	path := s.filePath(key)
	_, err := os.Stat(path)
	return err == nil
}

// Write stores the file content from reader to disk
func (s *Store) Write(key string, r io.Reader) error {
	path := s.filePath(key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

// Read opens the file and returns an io.ReadCloser
func (s *Store) Read(key string) (io.ReadCloser, error) {
	path := s.filePath(key)
	return os.Open(path)
}

// Delete removes the file from disk
func (s *Store) Delete(key string) error {
	path := s.filePath(key)
	return os.Remove(path)
}

// filePath generates the full path based on key using path hashing
func (s *Store) filePath(key string) string {
	// Example: abcd1234 → ./data/ab/cd/abcd1234
	sub1 := key[0:2]
	sub2 := key[2:4]
	return filepath.Join(s.BaseDir, sub1, sub2, key)
}
