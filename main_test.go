package main

import (
	"os"
	"testing"
)

func TestEncryptDecryptFile(t *testing.T) {
	// Create a temporary input file
	inputFile, err := os.CreateTemp("", "input.txt")
	if err != nil {
		t.Fatalf("failed to create temp input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Write some content to the input file
	content := []byte("This is a test content.")
	if _, err := inputFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp input file: %v", err)
	}
	inputFile.Close()

	// Create a temporary output file for encryption
	encryptedFile, err := os.CreateTemp("", "encrypted.txt")
	if err != nil {
		t.Fatalf("failed to create temp encrypted file: %v", err)
	}
	encryptedFile.Close()
	defer os.Remove(encryptedFile.Name())

	// Encrypt the file
	password := "testpassword"
	err = encryptFile(inputFile.Name(), encryptedFile.Name(), password)
	if err != nil {
		t.Fatalf("failed to encrypt file: %v", err)
	}

	// Create a temporary output file for decryption
	decryptedFile, err := os.CreateTemp("", "decrypted.txt")
	if err != nil {
		t.Fatalf("failed to create temp decrypted file: %v", err)
	}
	decryptedFile.Close()
	defer os.Remove(decryptedFile.Name())

	// Decrypt the file
	err = decryptFile(encryptedFile.Name(), decryptedFile.Name(), password)
	if err != nil {
		t.Fatalf("failed to decrypt file: %v", err)
	}

	// Read the decrypted content
	decryptedContent, err := os.ReadFile(decryptedFile.Name())
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	// Compare the decrypted content with the original content
	if string(decryptedContent) != string(content) {
		t.Fatalf("decrypted content does not match original content")
	}
}

func TestEncryptFileWithInvalidInput(t *testing.T) {
	// Try to encrypt a non-existent file
	err := encryptFile("nonexistent.txt", "output.txt", "password")
	if err == nil {
		t.Fatalf("expected error when encrypting non-existent file, got nil")
	}
}

func TestDecryptFileWithInvalidInput(t *testing.T) {
	// Try to decrypt a non-existent file
	err := decryptFile("nonexistent.txt", "output.txt", "password")
	if err == nil {
		t.Fatalf("expected error when decrypting non-existent file, got nil")
	}
}

func TestDecryptFileWithWrongPassword(t *testing.T) {
	// Create a temporary input file
	inputFile, err := os.CreateTemp("", "input.txt")
	if err != nil {
		t.Fatalf("failed to create temp input file: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Write some content to the input file
	content := []byte("This is a test content.")
	if _, err := inputFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp input file: %v", err)
	}
	inputFile.Close()

	// Create a temporary output file for encryption
	encryptedFile, err := os.CreateTemp("", "encrypted.txt")
	if err != nil {
		t.Fatalf("failed to create temp encrypted file: %v", err)
	}
	encryptedFile.Close()
	defer os.Remove(encryptedFile.Name())

	// Encrypt the file
	password := "testpassword"
	err = encryptFile(inputFile.Name(), encryptedFile.Name(), password)
	if err != nil {
		t.Fatalf("failed to encrypt file: %v", err)
	}

	// Create a temporary output file for decryption
	decryptedFile, err := os.CreateTemp("", "decrypted.txt")
	if err != nil {
		t.Fatalf("failed to create temp decrypted file: %v", err)
	}
	decryptedFile.Close()
	defer os.Remove(decryptedFile.Name())

	// Try to decrypt the file with the wrong password
	wrongPassword := "wrongpassword"
	err = decryptFile(encryptedFile.Name(), decryptedFile.Name(), wrongPassword)
	if err == nil {
		t.Fatalf("expected error when decrypting with wrong password, got nil")
	}
}
