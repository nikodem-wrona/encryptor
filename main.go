package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func encryptFile(inputFile, outputFile, password string) error {
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

func decryptFile(inputFile, outputFile, password string) error {
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

func main() {
	// Define command-line flags
	encryptFlag := flag.Bool("encrypt", false, "Encrypt the file")
	decryptFlag := flag.Bool("decrypt", false, "Decrypt the file")
	inputFileFlag := flag.String("input", "", "Path to the input file")
	outputFileFlag := flag.String("output", "", "Path to the output file")

	// Parse command-line flags
	flag.Parse()

	// Check if input and output file paths are provided
	if *inputFileFlag == "" || *outputFileFlag == "" {
		fmt.Println("Please specify both input and output file paths using --input and --output flags")
		return
	}

	// Read password from stdin
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	password = password[:len(password)-1] // Remove the newline character

	if *encryptFlag {
		err := encryptFile(*inputFileFlag, *outputFileFlag, password)
		if err != nil {
			fmt.Println("Error encrypting file:", err)
		}
	} else if *decryptFlag {
		err := decryptFile(*inputFileFlag, *outputFileFlag, password)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
		}
	} else {
		fmt.Println("Please specify either --encrypt or --decrypt")
	}
}
