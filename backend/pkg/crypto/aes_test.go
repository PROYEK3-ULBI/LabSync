package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptEmbedding(t *testing.T) {
	key := []byte("01234567890123456789012345678901") // 32 bytes

	original := []float64{0.12345, -0.6789, 0.0, 1.23456}

	ciphertext, err := EncryptEmbedding(key, original)
	if err != nil {
		t.Fatalf("EncryptEmbedding error: %v", err)
	}

	if bytes.Contains(ciphertext, []byte("0.12345")) {
		t.Errorf("ciphertext leaked plaintext numbers")
	}

	decrypted, err := DecryptEmbedding(key, ciphertext)
	if err != nil {
		t.Fatalf("DecryptEmbedding error: %v", err)
	}

	if len(decrypted) != len(original) {
		t.Fatalf("expected len %d, got %d", len(original), len(decrypted))
	}

	for i := range original {
		if decrypted[i] != original[i] {
			t.Errorf("index %d: expected %f, got %f", i, original[i], decrypted[i])
		}
	}
}

func TestDecryptWithInvalidKeyOrTamperedData(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	wrongKey := []byte("11234567890123456789012345678901")

	original := []float64{0.5, 0.5}
	ciphertext, _ := EncryptEmbedding(key, original)

	// Test wrong key
	_, err := DecryptEmbedding(wrongKey, ciphertext)
	if err == nil {
		t.Errorf("expected error when decrypting with wrong key")
	}

	// Test tampered ciphertext
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)-1] ^= 0xFF

	_, err2 := DecryptEmbedding(key, tampered)
	if err2 == nil {
		t.Errorf("expected error when decrypting tampered data")
	}
}
