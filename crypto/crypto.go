package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
)

var (
	ErrInvalidKeySize = errors.New("invalid AES key size (must be 16, 24, or 32 bytes)")
)

func Encrypt(key, iv, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aes.BlockSize {
		return nil, errors.New("invalid IV size")
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}

func Decrypt(key, iv, ciphertext []byte) ([]byte, error) {
	return Encrypt(key, iv, ciphertext) // CTR模式对称
}

func EncryptStream(key, iv []byte, in io.Reader, out io.Writer) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: out}
	_, err = io.Copy(writer, in)
	return err
}

func DecryptStream(key, iv []byte, in io.Reader, out io.Writer) error {
	return EncryptStream(key, iv, in, out)
}
