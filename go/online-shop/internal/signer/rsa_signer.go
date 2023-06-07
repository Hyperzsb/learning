package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

type RSASigner struct {
	pk *rsa.PublicKey
	sk *rsa.PrivateKey
}

func NewRSASigner(pk *rsa.PublicKey, sk *rsa.PrivateKey) *RSASigner {
	return &RSASigner{
		pk: pk,
		sk: sk,
	}
}

func (s *RSASigner) Sign(data string) (string, error) {
	dataHash := sha256.Sum256([]byte(data))

	timeString := time.Now().Format(time.UnixDate)
	timeEncoded := base64.StdEncoding.EncodeToString([]byte(timeString))

	dataTime := make([]byte, len(dataHash)+len(timeEncoded))

	idx := 0
	for i := range dataHash {
		dataTime[idx] = dataHash[i]
		idx++
	}
	for i := range timeEncoded {
		dataTime[idx] = timeEncoded[i]
		idx++
	}

	dataTimeHash := sha256.Sum256(dataTime)
	dataTimeSignature, err := rsa.SignPKCS1v15(rand.Reader, s.sk, crypto.SHA256, dataTimeHash[:])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", timeEncoded, hex.EncodeToString(dataTimeSignature)), nil
}

func (s *RSASigner) Verify(data, signature string) (time.Time, error) {
	dataHash := sha256.Sum256([]byte(data))

	signatureParts := strings.Split(signature, ".")
	if len(signatureParts) != 2 {
		return time.Time{}, errors.New("incomplete signature")
	}

	timeEncoded := signatureParts[0]

	dataTimeSignature, err := hex.DecodeString(signatureParts[1])
	if err != nil {
		return time.Time{}, err
	}

	dataTime := make([]byte, len(dataHash)+len(timeEncoded))

	idx := 0
	for i := range dataHash {
		dataTime[idx] = dataHash[i]
		idx++
	}
	for i := range timeEncoded {
		dataTime[idx] = timeEncoded[i]
		idx++
	}

	dataTimeHash := sha256.Sum256(dataTime)

	err = rsa.VerifyPKCS1v15(s.pk, crypto.SHA256, dataTimeHash[:], dataTimeSignature)
	if err != nil {
		return time.Time{}, err
	}

	timeString, err := base64.StdEncoding.DecodeString(timeEncoded)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.Parse(time.UnixDate, string(timeString))
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
