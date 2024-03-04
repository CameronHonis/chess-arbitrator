package auth

import (
	"crypto/sha256"
	"encoding/hex"
	. "github.com/CameronHonis/chess-arbitrator/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func CreateClient(conn *websocket.Conn, cleanup func(*Client)) *Client {
	pubKey, priKey := GenerateKeyset()

	return NewClient(pubKey, priKey, conn, cleanup)
}

func GenerateKeyset() (publicKey Key, privateKey Key) {
	priKey := GeneratePriKey()
	pubKey := sha256.Sum256([]byte(priKey))
	return Key(hex.EncodeToString(pubKey[:])), Key(priKey)
}

func GeneratePriKey() Key {
	return Key(uuid.New().String())
}

func ValidatePrivateKey(publicKey Key, privateKey Key) bool {
	publicKeyFromPrivateKey := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(publicKeyFromPrivateKey[:]) == string(publicKey)
}
