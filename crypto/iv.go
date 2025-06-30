package crypto

import (
	"crypto/rand"
	"errors"
)

func NewIV() ([]byte, error) {
	iv := make([]byte, 16)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, errors.New("failed to generate IV")
	}
	return iv, nil
}
