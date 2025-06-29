package p2p

import (
	"bytes"
	"encoding/gob"
)

// Encode serializes a struct to []byte using GOB
func Encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode deserializes a GOB []byte into a struct pointer
func Decode[T any](b []byte) (*T, error) {
	var result T
	buf := bytes.NewReader(b)
	if err := gob.NewDecoder(buf).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
