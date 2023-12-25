package server

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

func GenerateKeyset() (publicKey Key, privateKey Key) {
	privKey := uuid.New().String()
	pubKey := sha256.Sum256([]byte(privKey))
	return Key(hex.EncodeToString(pubKey[:])), Key(privKey)
}

func ValidatePrivateKey(publicKey Key, privateKey Key) bool {
	publicKeyFromPrivateKey := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(publicKeyFromPrivateKey[:]) == string(publicKey)
}
