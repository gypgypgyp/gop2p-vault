package main

import (
	"fmt"
	// "os"
	"strings"

	"gop2p-vault/store"
)

func main() {
	// Create a new store under ./data directory
	s := store.New("./data")

	// Simulated input
	content := "Hello P2P Vault Storage!"
	reader := strings.NewReader(content)

	// Generate key by hashing content
	key, err := store.HashKey(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	fmt.Println("[Generated Key]:", key)

	// Write to store
	if err := s.Write(key, reader); err != nil {
		panic(err)
	}
	fmt.Println("[Write]: Success")

	// Check if file exists
	exists := s.Has(key)
	fmt.Println("[Has]:", exists)

	// Read and print content
	r, err := s.Read(key)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	data := make([]byte, len(content))
	r.Read(data)
	fmt.Println("[Read]:", string(data))

	// Delete file
	if err := s.Delete(key); err != nil {
		panic(err)
	}
	fmt.Println("[Delete]: Success")

	// Verify deletion
	fmt.Println("[Has after delete]:", s.Has(key))
}
