package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"os"
)

var authManager *AuthManager

type AuthManager struct {
	chessBotKey string
}

func GetAuthManager() *AuthManager {
	if authManager == nil {
		authManager = &AuthManager{}
	}
	return authManager
}

func (am *AuthManager) UpgradeAuth(clientKey string, secret string) (string, error) {
	chessBotKey, ok := os.LookupEnv("BOT_CLIENT_SECRET")
	if !ok {
		return "", fmt.Errorf("could not determine chess bot secret")
	}

	switch secret {
	case chessBotKey:
		if am.chessBotKey != "" {
			return "", fmt.Errorf("chess bot already authenticated")
		}
		am.chessBotKey = clientKey
		return "chessBot", nil
	default:
		return "", fmt.Errorf("unrecognized secret")
	}
}

func GenerateKeyset() (string, string) {
	privateKey := uuid.New().String()
	publicKey := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(publicKey[:]), privateKey
}

func ValidatePrivateKey(publicKey string, privateKey string) bool {
	publicKeyFromPrivateKey := sha256.Sum256([]byte(privateKey))
	return hex.EncodeToString(publicKeyFromPrivateKey[:]) == publicKey
}

func (am *AuthManager) ValidateAuthInMessage(msg *Message) error {
	isValidAuth := ValidatePrivateKey(msg.SenderKey, msg.PrivateKey)
	if !isValidAuth {
		return fmt.Errorf("invalid auth")
	}
	return nil
}
