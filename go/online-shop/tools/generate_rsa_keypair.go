package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	bits      = 2048
	plainText = "This is a plain text."
)

// demo demonstrates basic usages of generating an RSA key pair, encrypting
// and decrypting data with the key pair, signing and validating data with
// the key pair, converting and saving the key pair in X.509 format.
func demo() error {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	pk, sk := &key.PublicKey, key

	// Encryption and decryption

	fmt.Printf("Plain Text: %s\n\n", plainText)

	encryptedText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, []byte(plainText), nil)
	if err != nil {
		return err
	}

	fmt.Printf("Encrypted Text: %x\n\n", encryptedText)

	decryptedText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, sk, encryptedText, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Decrypted Text: %s\n\n", decryptedText)

	// Signature and validation

	hash := sha256.Sum256([]byte(plainText))
	signature, err := rsa.SignPKCS1v15(rand.Reader, sk, crypto.SHA256, hash[:])

	fmt.Printf("Signature: %x\n\n", string(signature))

	err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, hash[:], signature)
	if err != nil {
		fmt.Printf("Signature is invalid\n\n")
	} else {
		fmt.Printf("Signature is valid\n\n")
	}

	signature[0] = 0
	err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, hash[:], signature)
	if err != nil {
		fmt.Printf("Signature is invalid\n\n")
	} else {
		fmt.Printf("Signature is valid\n\n")
	}

	// X.509 conversion

	pkBytes := x509.MarshalPKCS1PublicKey(pk)
	pkPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pkBytes,
	})
	fmt.Println(string(pkPem))

	skBytes := x509.MarshalPKCS1PrivateKey(sk)
	skPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: skBytes,
	})
	fmt.Println(string(skPem))

	loadedPKBlock, _ := pem.Decode(pkPem)
	if loadedPKBlock.Bytes == nil || loadedPKBlock.Type != "RSA PUBLIC KEY" {
		return errors.New("invalid public key raw data")
	}

	loadedPK, err := x509.ParsePKCS1PublicKey(loadedPKBlock.Bytes)
	if err != nil {
		return err
	}

	loadedSKBlock, _ := pem.Decode(skPem)
	if loadedSKBlock.Bytes == nil || loadedSKBlock.Type != "RSA PRIVATE KEY" {
		return errors.New("invalid private key raw data")
	}

	loadedSK, err := x509.ParsePKCS1PrivateKey(loadedSKBlock.Bytes)
	if err != nil {
		return err
	}

	// Encryption and decryption again

	fmt.Printf("Plain Text: %s\n\n", plainText)

	encryptedText, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, loadedPK, []byte(plainText), nil)
	if err != nil {
		return err
	}

	fmt.Printf("Encrypted Text: %x\n\n", encryptedText)

	decryptedText, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, loadedSK, encryptedText, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Decrypted Text: %s\n\n", decryptedText)

	return nil
}

func generate(path string) error {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	pk, sk := &key.PublicKey, key

	// Save the public key

	pkBytes := x509.MarshalPKCS1PublicKey(pk)
	pkPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pkBytes,
	})

	pkFile, err := os.Create(fmt.Sprintf("%s/rsa.pub", path))
	if err != nil {
		return err
	}
	defer pkFile.Close()

	_, err = pkFile.Write(pkPem)
	if err != nil {
		return err
	}

	// Save the private key

	skBytes := x509.MarshalPKCS1PrivateKey(sk)
	skPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: skBytes,
	})

	skFile, err := os.Create(fmt.Sprintf("%s/rsa.pem", path))
	if err != nil {
		return err
	}
	defer skFile.Close()

	_, err = skFile.Write(skPem)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := generate(".data"); err != nil {
		log.Fatal(err)
	}
}
