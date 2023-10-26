package server

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

func GenerateKeyset() (string, string) {
	privateKey := uuid.New().String()
	publicKey := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(publicKey[:]), privateKey
}

func ValidatePrivateKey(publicKey string, privateKey string) bool {
	publicKeyFromPrivateKey := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(publicKeyFromPrivateKey[:]) == publicKey
}
