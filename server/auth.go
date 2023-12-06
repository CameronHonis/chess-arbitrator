package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
)

type AuthManagerI interface {
	UpgradeAuth(clientKey string, secret string) error
	ValidateAuthInMessage(msg *Message) error
}

var authManager *AuthManager

type AuthManager struct {
	userClientsManager UserClientsManagerI

	chessBotKey string
	mu          sync.Mutex
}

func GetAuthManager() *AuthManager {
	if authManager == nil {
		authManager = &AuthManager{
			userClientsManager: GetUserClientsManager(),
		}
	}
	return authManager
}

func (am *AuthManager) UpgradeAuth(clientKey string, secret string) error {
	chessBotKey, ok := os.LookupEnv("BOT_CLIENT_SECRET")
	if !ok {
		return fmt.Errorf("could not determine chess bot secret")
	}

	switch secret {
	case chessBotKey:
		if am.chessBotKey != "" {
			return fmt.Errorf("chess bot already authenticated")
		}
		am.chessBotKey = clientKey
		msg := Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_UPGRADE_AUTH_GRANTED,
			Content: &UpgradeAuthGrantedMessageContent{
				UpgradedToRole: "chessBot",
			},
		}
		return GetUserClientsManager().DirectMessage(&msg, clientKey)
	default:
		msg := Message{
			Topic:       "directMessage",
			ContentType: CONTENT_TYPE_UPGRADE_AUTH_DENIED,
			Content: &UpgradeAuthDeniedMessageContent{
				Reason: "unrecognized secret",
			},
		}
		return GetUserClientsManager().DirectMessage(&msg, clientKey)
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
