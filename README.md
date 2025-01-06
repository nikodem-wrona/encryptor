# Encryptor

Encryptor is a command-line tool for encrypting and decrypting files using AES encryption.

### Usage
To use the program, simply run it, enter the encryption/decryption password when prompted, and select the file you wish to encrypt or decrypt.


### Password
You will be prompted to enter a password when encrypting or decrypting a file. The password is used to derive the encryption key.

## Running Tests

To run the tests, use the following command:

```sh
go test ./...
```

## Build

To build the program, use the following command:

```sh
go build -o encryptor main.go
```

This will create an executable named `encryptor` in the current directory.

## Note
The source file will be deleted after encryption.
Ensure you remember the password used for encryption, as it is required for decryption.

## License
This project is licensed under the MIT License
