package protocol

import (
	"bytes"
	"encoding/gob"
)

type Metadata struct {
	Name     string
	Size     int64
	Hash     string
	MimeType string
}

func EncodeMetadata(md *Metadata) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(md)
	return buf.Bytes(), err
}

func DecodeMetadata(data []byte) (*Metadata, error) {
	var md Metadata
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(&md)
	return &md, err
}
