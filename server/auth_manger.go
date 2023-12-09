package server

import (
	"fmt"
	"os"
	"sync"
)

type AuthManagerI interface {
	GetBotKey() (string, error)
	//GetRole() (string, error)
	UpgradeAuth(clientKey string, secret string) (role string, err error)
	ValidateAuthInMessage(msg *Message) error
	ValidateClientForTopic(clientKey string, topic MessageTopic) error
}

var authManager *AuthManager

type AuthManager struct {
	userClientsManager UserClientsManagerI

	chessBotKey string

	mu sync.Mutex
}

func GetAuthManager() *AuthManager {
	if authManager == nil {
		authManager = &AuthManager{
			userClientsManager: GetUserClientsManager(),
		}
	}
	return authManager
}

func (am *AuthManager) GetBotKey() (string, error) {
	if am.chessBotKey == "" {
		return "", fmt.Errorf("chess bot not authenticated")
	}
	return am.chessBotKey, nil
}

func (am *AuthManager) setBotKey(botKey string) {
	am.chessBotKey = botKey
}

func (am *AuthManager) UpgradeAuth(clientKey string, secret string) (role string, err error) {
	chessBotKey, ok := os.LookupEnv("BOT_CLIENT_SECRET")
	if !ok {
		return "", fmt.Errorf("could not determine chess bot secret")
	}

	switch secret {
	case chessBotKey:
		if am.chessBotKey != "" {
			return "", fmt.Errorf("chess bot already authenticated")
		}

		am.setBotKey(clientKey)
		return "chessBot", nil
	default:
		return "", fmt.Errorf("invalid secret")
	}

}

func (am *AuthManager) ValidateAuthInMessage(msg *Message) error {
	isValidAuth := ValidatePrivateKey(msg.SenderKey, msg.PrivateKey)
	if !isValidAuth {
		return fmt.Errorf("invalid auth")
	}
	return nil
}

func (am *AuthManager) ValidateClientForTopic(clientKey string, topic MessageTopic) error {
	//role, getRoleErr := am.userClientsManager.GetRole(clientKey)
	return nil
}
