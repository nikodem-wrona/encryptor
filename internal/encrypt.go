package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/xdg-go/pbkdf2"
)

func EncryptFile(inputFile, outputFile, password string) error {
	// Open the input file
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive a key from the password using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Generate a random nonce (IV)
	nonce := make([]byte, 12) // GCM nonce is 12 bytes
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

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

	// Encrypt the file content
	ciphertext := aesGCM.Seal(nil, nonce, input, nil)

	// Write the salt, nonce, and ciphertext to the output file
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	if _, err := output.Write(salt); err != nil {
		return fmt.Errorf("failed to write salt: %w", err)
	}
	if _, err := output.Write(nonce); err != nil {
		return fmt.Errorf("failed to write nonce: %w", err)
	}
	if _, err := output.Write(ciphertext); err != nil {
		return fmt.Errorf("failed to write ciphertext: %w", err)
	}

	// Delete the source file after encryption
	err = os.Remove(inputFile)
	if err != nil {
		return fmt.Errorf("failed to delete source file: %w", err)
	}

	return nil
}
