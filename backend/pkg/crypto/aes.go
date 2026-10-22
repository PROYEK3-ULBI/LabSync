package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// EncryptEmbedding mengenkripsi slice float64 embedding ke ciphertext BYTEA menggunakan AES-256-GCM
func EncryptEmbedding(key []byte, embedding []float64) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("kunci AES-256 harus tepat 32 bytes")
	}

	plaintext, err := json.Marshal(embedding)
	if err != nil {
		return nil, fmt.Errorf("gagal serialisasi embedding: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptEmbedding mendekripsi ciphertext BYTEA kembali menjadi slice float64 embedding
func DecryptEmbedding(key []byte, ciphertext []byte) ([]float64, error) {
	if len(key) != 32 {
		return nil, errors.New("kunci AES-256 harus tepat 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("panjang ciphertext terlalu pendek")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, errors.New("dekripsi ciphertext biometrik gagal atau integritas data tidak valid")
	}

	var embedding []float64
	if err := json.Unmarshal(plaintext, &embedding); err != nil {
		return nil, fmt.Errorf("gagal deserialisasi embedding: %w", err)
	}

	return embedding, nil
}
