package server

import (
	"fmt"
	"io"
	"os"
	"time"
	"bytes"

	"gop2p-vault/p2p"
	"gop2p-vault/store"
	"gop2p-vault/crypto"
)

var secretKey = []byte("gop2p-vault-key1") // 16-byte key (AES-128)

// HandleUpload processes an upload message and stores the file locally
func HandleUpload(data []byte) (string, error) {
	iv, err := crypto.NewIV()
	if err != nil {
		return "", fmt.Errorf("iv gen failed: %w", err)
	}

	encData, err := crypto.Encrypt(secretKey, iv, data)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	combined := append(iv, encData...) // prepend IV

	key, err := store.HashKeyBytes(combined)
	if err != nil {
		return "", fmt.Errorf("failed to compute file hash: %w", err)
	}

	s := store.New("./data")
	err = s.Write(key, bytes.NewReader(combined))
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return key, nil
}

// HandleDownload loads the requested file content by key and returns a Message
func HandleDownload(fileKey string) (*p2p.Message, error) {
	s := store.New("./data")
	reader, err := s.Read(fileKey)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	fullData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	if len(fullData) < 16 {
		return nil, fmt.Errorf("data too short")
	}
	iv := fullData[:16]
	ciphertext := fullData[16:]

	plaintext, err := crypto.Decrypt(secretKey, iv, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return &p2p.Message{
		Type: "download_result",
		Data: plaintext,
	}, nil
}

// HandleDownloadResult saves the downloaded content to a local temp file
func HandleDownloadResult(data []byte) (string, error) {
	tmpPath := fmt.Sprintf("./data/downloaded_%d", time.Now().UnixNano())
	err := os.WriteFile(tmpPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	return tmpPath, nil
}
