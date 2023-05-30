package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type Token struct {
	UserID     int       `json:"-"`
	Scope      string    `json:"-"`
	Text       []byte    `json:"text"`
	Hash       [32]byte  `json:"-"`
	ExpiryTime time.Time `json:"expiry_time"`
}

func NewToken(userID int, scope string, ttl time.Duration) (Token, error) {
	token := Token{
		UserID:     userID,
		Scope:      scope,
		ExpiryTime: time.Now().Add(ttl),
	}

	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return token, err
	}

	token.Text = make([]byte, base64.StdEncoding.EncodedLen(len(nonce)))
	base64.StdEncoding.WithPadding(base64.NoPadding).Encode(token.Text, nonce)
	token.Hash = sha256.Sum256(token.Text)
	return token, nil
}
