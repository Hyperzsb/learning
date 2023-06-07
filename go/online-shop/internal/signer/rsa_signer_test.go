package signer

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

const (
	data = "This is a piece of data."
)

var (
	pk *rsa.PublicKey
	sk *rsa.PrivateKey
)

func setup() error {
	// Inject the file path of crypto keys via environment variables
	rsaPKPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	if rsaPKPath == "" {
		rsaPKPath = ".config/rsa.pub"
	}
	rsaSKPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	if rsaSKPath == "" {
		rsaSKPath = ".config/rsa.pem"
	}

	rsaPKBytes, err := os.ReadFile(rsaPKPath)
	if err != nil {
		return err
	}

	rsaSKBytes, err := os.ReadFile(rsaSKPath)
	if err != nil {
		return err
	}

	rsaPKBlock, _ := pem.Decode(rsaPKBytes)
	rsaSKBlock, _ := pem.Decode(rsaSKBytes)

	rsaPK, err := x509.ParsePKCS1PublicKey(rsaPKBlock.Bytes)
	if err != nil {
		return err
	}

	rsaSK, err := x509.ParsePKCS1PrivateKey(rsaSKBlock.Bytes)
	if err != nil {
		return err
	}

	pk = rsaPK
	sk = rsaSK

	return nil
}

func TestRSASigner_Sign(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatal("Unable to load keys")
	}

	signer := NewRSASigner(pk, sk)
	signature, err := signer.Sign(data)
	if err != nil {
		t.Error("Unable to sign the data")
	}

	t.Logf("Data.Signature: %s.%s", data, signature)
}

func TestRSASigner_Verify(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatal("Unable to load keys")
	}

	signer := NewRSASigner(pk, sk)
	signature, err := signer.Sign(data)
	if err != nil {
		t.Error("Unable to sign the data")
	}

	tm, err := signer.Verify(data, signature)
	if err != nil {
		t.Error("Unable to verify the signature")
	}

	t.Logf("Signature time: %v", tm)
}
