package server

import (
	"fmt"
	"io"
	"os"
	"time"
	"bytes"

	"gop2p-vault/p2p"
	"gop2p-vault/store"
)

// HandleUpload processes an upload message and stores the file locally
func HandleUpload(data []byte) (string, error) {
	key, err := store.HashKeyBytes(data)
	if err != nil {
		return "", fmt.Errorf("failed to compute file hash: %w", err)
	}

	s := store.New("./data")
	// err = s.Write(key, store.BytesReader(data))
	err = s.Write(key, bytes.NewReader(data))
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
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	return &p2p.Message{
		Type: "download_result",
		Data: content,
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
