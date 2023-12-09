package server

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

func GenerateKeyset() (publicKey string, privateKey string) {
	privKey := uuid.New().String()
	pubKey := sha256.Sum256([]byte(privKey))
	return hex.EncodeToString(pubKey[:]), privKey
}

func ValidatePrivateKey(publicKey string, privateKey string) bool {
	publicKeyFromPrivateKey := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(publicKeyFromPrivateKey[:]) == publicKey
}
