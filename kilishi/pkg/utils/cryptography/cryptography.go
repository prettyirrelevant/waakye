package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
)

// Encrypt takes a plaintext string, a secret key in hexadecimal format, and an
// initialization vector in hexadecimal format, and returns the encrypted
// ciphertext as a base64-encoded string using AES-256 in CBC mode.
func Encrypt(plaintext string, secretKeyHex string, initializationVectorHex string) (string, error) {
	// Decode the secret key from a hex string to a byte slice
	secretKey, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block with the decoded secret key
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	// Pad the plaintext message to the nearest multiple of the block size
	plaintextBytes := []byte(plaintext)
	paddedPlaintext := pad(plaintextBytes, block.BlockSize())

	// Create a byte slice to hold the encrypted ciphertext
	ciphertext := make([]byte, len(paddedPlaintext))

	// Decode the initialization vector from a hex string to a byte slice
	iv, err := hex.DecodeString(initializationVectorHex)
	if err != nil {
		return "", err
	}

	// Create a new CBC mode encrypter with the AES cipher block and initialization vector
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt the padded plaintext message and store the result in the ciphertext byte slice
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Encode the ciphertext byte slice as a base64-encoded string
	ciphertextStr := base64.StdEncoding.EncodeToString(ciphertext)

	return ciphertextStr, nil
}

// Decrypt decrypts the given ciphertext string using the provided secret key and initialization vector.
func Decrypt(ciphertextStr string, secretKeyHex string, initializationVectorHex string) (string, error) {
	// Decode the secret key from a hex string to a byte slice
	key, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	// Decode the base64-encoded ciphertext string to a byte slice
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a byte slice to hold the decrypted plaintext
	plaintext := make([]byte, len(ciphertext))

	// Decode the initialization vector from a hex string to a byte slice
	iv, err := hex.DecodeString(initializationVectorHex)
	if err != nil {
		return "", err
	}

	// Create a new CBC mode decrypter
	decrypter := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the ciphertext byte slice and store the result in the plaintext byte slice
	decrypter.CryptBlocks(plaintext, ciphertext)

	// Remove the PKCS7 padding from the plaintext byte slice
	unpaddedPlaintext := unpad(plaintext)

	// Convert the unpadded plaintext byte slice to a string and return it
	return string(unpaddedPlaintext), nil
}

// pad adds PKCS7 padding to the message byte slice, such that the length
// of the padded message is a multiple of the specified block size.
func pad(message []byte, blockSize int) []byte {
	padding := blockSize - len(message)%blockSize
	// the value of each byte in the padding is equal to the number of padding bytes added
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	// append the padding to the original message
	return append(message, padtext...)
}

// unpad removes PKCS7 padding from the message byte slice, returning the original
// message with the padding removed.
func unpad(message []byte) []byte {
	// get the value of the last byte in the message, which represents the padding
	padding := int(message[len(message)-1])
	// remove the padding from the message by returning a slice that excludes the padding bytes
	return message[:len(message)-padding]
}
