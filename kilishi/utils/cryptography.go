package utils

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
func Encrypt(plaintext, secretKeyHex, initializationVectorHex string) (string, error) {
	secretKey, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	paddedPlaintext := pad(plaintextBytes, block.BlockSize())

	ciphertext := make([]byte, len(paddedPlaintext))

	iv, err := hex.DecodeString(initializationVectorHex)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	ciphertextStr := base64.StdEncoding.EncodeToString(ciphertext)
	return ciphertextStr, nil
}

// Decrypt decrypts the given ciphertext string using the provided secret key and initialization vector.
func Decrypt(ciphertextStr, secretKeyHex, initializationVectorHex string) (string, error) {
	key, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertext))

	iv, err := hex.DecodeString(initializationVectorHex)
	if err != nil {
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(plaintext, ciphertext)

	unpaddedPlaintext := unpad(plaintext)
	return string(unpaddedPlaintext), nil
}

// pad adds PKCS7 padding to the message byte slice, such that the length
// of the padded message is a multiple of the specified block size.
func pad(message []byte, blockSize int) []byte {
	padding := blockSize - len(message)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(message, padtext...)
}

// unpad removes PKCS7 padding from the message byte slice, returning the original
// message with the padding removed.
func unpad(message []byte) []byte {
	padding := int(message[len(message)-1])
	return message[:len(message)-padding]
}
