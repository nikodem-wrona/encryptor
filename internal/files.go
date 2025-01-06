package internal

import (
	"fmt"
	"os"
)

func GetFilePathByFileName(fileName string) string {
	files, err := os.ReadDir(".")

	if err != nil {
		fmt.Println("Error reading directory")
		os.Exit(1)
	}

	for _, file := range files {
		if file.Name() == fileName {
			return file.Name()
		}
	}

	return ""
}

func GetFilePathForEncryptedFile(fileName string) string {
	return fileName + ".enc"
}

func IsEncryptedFile(fileName string) bool {
	return fileName[len(fileName)-4:] == ".enc"
}


func FindAllFiledInCurrentDir() []string {
	files, err := os.ReadDir(".")

	if err != nil {
		fmt.Println("Error reading directory")
		os.Exit(1)
	}

	var fileNames []string

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames
}