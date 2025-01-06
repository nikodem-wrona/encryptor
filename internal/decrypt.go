package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"github.com/xdg-go/pbkdf2"
)

func DecryptFile(inputFile, outputFile, password string) error {
	// Read the encrypted file
	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Ensure the file has enough data for salt (16 bytes) and nonce (12 bytes)
	if len(encryptedData) < 28 {
		return errors.New("file too short to contain salt, nonce, and ciphertext")
	}

	// Extract salt, nonce, and ciphertext
	salt := encryptedData[:16]
	nonce := encryptedData[16:28]
	ciphertext := encryptedData[28:]

	// Derive the key using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the ciphertext
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Write the plaintext to the output file
	err = os.WriteFile(outputFile, plaintext, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	return nil
}
