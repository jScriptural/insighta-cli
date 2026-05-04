package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"crypto/rand"
)

func GeneratePKCE() (string, string) {
	b := make([]byte, 32)
	rand.Read(b)
	verifier := base64.RawURLEncoding.EncodeToString(b)

	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return verifier, challenge
}

