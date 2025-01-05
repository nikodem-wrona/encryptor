# Encryptor

Encryptor is a command-line tool for encrypting and decrypting files using AES encryption.

## Usage

To use the Encryptor, you need to specify whether you want to encrypt or decrypt a file, along with the input and output file paths. You will also be prompted to enter a password.

### Encrypting a File

To encrypt a file, use the `--encrypt` flag along with the `--input` and `--output` flags to specify the input and output file paths.

```sh
go run main.go --encrypt --input <path_to_input_file> --output <path_to_output_file>
```

### Decrypting a File

To decrypt a file, use the --decrypt flag along with the --input and --output flags to specify the input and output file paths.

```sh
go run main.go --decrypt --input <path_to_input_file> --output <path_to_output_file>
```

### Password
You will be prompted to enter a password when encrypting or decrypting a file. The password is used to derive the encryption key.

### Flags
- --encrypt: Encrypt the file.
- --decrypt: Decrypt the file.
- --input: Path to the input file.
- --output: Path to the output file.

## Note
The source file will be deleted after encryption.
Ensure you remember the password used for encryption, as it is required for decryption.

## License
This project is licensed under the MIT License.