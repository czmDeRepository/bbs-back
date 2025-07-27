package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"regexp"
	"testing"
)

func TestDecryptRSA(t *testing.T) {
	// First, we need to create a valid encrypted string to test decryption
	// Parse the private key
	block, _ := pem.Decode(RsaPrivateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	// Get the public key
	publicKey := &privateKey.PublicKey

	// Test data
	plainText := []byte("test_password")

	// Encrypt with public key
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		t.Fatalf("Failed to encrypt test data: %v", err)
	}

	// Encode to base64
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)

	// Test decryption
	decrypted, err := DecryptRSA(encodedCipherText)
	if err != nil {
		t.Errorf("DecryptRSA() error = %v", err)
		return
	}

	// Check if decryption was successful
	if string(decrypted) != string(plainText) {
		t.Errorf("DecryptRSA() = %v, want %v", string(decrypted), string(plainText))
	}

	// Test with invalid input
	_, err = DecryptRSA("invalid_base64_string")
	if err == nil {
		t.Errorf("DecryptRSA() with invalid input should return error")
	}
}

func TestDecryptPassword(t *testing.T) {
	// First, we need to create a valid encrypted string to test decryption
	// Parse the private key
	block, _ := pem.Decode(RsaPrivateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	// Get the public key
	publicKey := &privateKey.PublicKey

	// Test data
	plainText := []byte("test_password")

	// Encrypt with public key
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		t.Fatalf("Failed to encrypt test data: %v", err)
	}

	// Encode to base64
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)

	// Test DecryptPassword
	decryptedPassword, err := DecryptPassword(encodedCipherText)
	if err != nil {
		t.Errorf("DecryptPassword() error = %v", err)
		return
	}

	// Check if the result is a valid MD5 hash (32 hex characters)
	if match, _ := regexp.MatchString("^[0-9a-f]{32}$", decryptedPassword); !match {
		t.Errorf("DecryptPassword() = %v, not a valid MD5 hash", decryptedPassword)
	}

	// Test with invalid input
	_, err = DecryptPassword("invalid_base64_string")
	if err == nil {
		t.Errorf("DecryptPassword() with invalid input should return error")
	}
}