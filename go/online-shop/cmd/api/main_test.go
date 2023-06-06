package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

// TestConfiguration_loadCrypto tests whether the loadCrypto function can load
// crypto keys from files, and whether we can use these keys to perform encryption,
// decryption, signature, and verification operations.
// To use this test case, you need to inject the corresponding environment variables
// indicating the locations of the key files, e.g.:
// $ env RSA_PUBLIC_KEY_PATH=".config/rsa.pub" RSA_PRIVATE_KEY_PATH=".config/rsa.pem" go test .
func TestConfiguration_loadCrypto(t *testing.T) {
	config := configuration{}

	if err := config.loadCrypto(); err != nil {
		t.Fatalf("Unable to load crypto: %v", err)
	}

	plainText := "This is a plain text."

	// Test encryption and decryption
	encryptedText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, config.crypto.rsa.pk, []byte(plainText), nil)
	if err != nil {
		t.Errorf("Unable to encrypt the plain text: %v", err)
	}

	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, config.crypto.rsa.sk, encryptedText, nil)
	if err != nil {
		t.Errorf("Unable to decrypt the plain text: %v", err)
	}

	if string(decryptedText) != plainText {
		t.Errorf("Encryption and decryption failed: %v", err)
	}

	// Test signature and verification
	hash := sha256.Sum256([]byte(plainText))

	signature, err := rsa.SignPKCS1v15(rand.Reader, config.crypto.rsa.sk, crypto.SHA256, hash[:])
	if err != nil {
		t.Error("Unable to sign the plain text")
	}

	err = rsa.VerifyPKCS1v15(config.crypto.rsa.pk, crypto.SHA256, hash[:], signature)
	if err != nil {
		t.Error("Unable to verify the signature")
	}
}
