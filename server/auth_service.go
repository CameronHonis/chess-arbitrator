package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"os"
	"sync"
)

type RoleName string

const (
	PLEB RoleName = "PLEB"
	BOT           = "BOT"
)

var ENV_NAME_BY_ROLE_NAME = map[RoleName]string{
	PLEB: "PLEB_SECRET",
	BOT:  "BOT_SECRET",
}

type AuthConfig struct {
	ConfigI
}

type AuthenticationServiceI interface {
	GetRole(clientKey string) (RoleName, error)

	AddClient(clientKey string)
	UpgradeAuth(clientKey string, roleName RoleName, secret string) error
	RemoveClient(clientKey string) error

	ValidateSecret(roleName RoleName, secret string) error
	ValidateAuthInMessage(msg *Message) error
	ValidateClientForTopic(clientKey string, topic MessageTopic) error

	getSecret(role RoleName) (string, error)
	setRole(clientKey string, role RoleName) error
}

type AuthenticationService struct {
	Service[*AuthConfig]

	__dependencies__ Marker
	LoggerService    LoggerServiceI

	__state__    Marker
	roleByClient map[string]RoleName
	mu           sync.Mutex
}

func NewAuthenticationService(config *AuthConfig) *AuthenticationService {
	authService := &AuthenticationService{}
	authService.Service = *NewService(authService, config)
	return authService
}

func (am *AuthenticationService) GetRole(clientKey string) (RoleName, error) {
	am.mu.Lock()
	role, ok := am.roleByClient[clientKey]
	am.mu.Unlock()
	if !ok {
		return "", fmt.Errorf("could not find role for client %s", clientKey)
	}
	return role, nil
}

func (am *AuthenticationService) AddClient(clientKey string) {
	am.mu.Lock()
	am.roleByClient[clientKey] = PLEB
	am.mu.Unlock()
}

func (am *AuthenticationService) UpgradeAuth(clientKey string, roleName RoleName, secret string) error {
	validSecretErr := am.ValidateSecret(roleName, secret)
	if validSecretErr != nil {
		return validSecretErr
	}
	return am.setRole(clientKey, roleName)
}

func (am *AuthenticationService) RemoveClient(clientKey string) error {
	_, roleErr := am.GetRole(clientKey)
	if roleErr != nil {
		am.LoggerService.LogRed(ENV_SERVER, fmt.Sprintf("attempt to remove non-existant client %s", clientKey))
		return nil
	}
	am.mu.Lock()
	delete(am.roleByClient, clientKey)
	am.mu.Unlock()
	return nil
}

func (am *AuthenticationService) ValidateSecret(roleName RoleName, secret string) error {
	expSecret, secretErr := am.getSecret(roleName)
	if secretErr != nil {
		return secretErr
	}
	if secret != expSecret {
		return fmt.Errorf("invalid secret")
	}
	return nil
}

func (am *AuthenticationService) ValidateAuthInMessage(msg *Message) error {
	isValidAuth := ValidatePrivateKey(msg.SenderKey, msg.PrivateKey)
	if !isValidAuth {
		return fmt.Errorf("invalid auth")
	}
	return nil
}

func (am *AuthenticationService) ValidateClientForTopic(clientKey string, topic MessageTopic) error {
	return nil
}

func (am *AuthenticationService) setRole(clientKey string, role RoleName) error {
	_, getRoleErr := am.GetRole(clientKey)
	if getRoleErr != nil {
		return getRoleErr
	}
	am.mu.Lock()
	am.roleByClient[clientKey] = role
	am.mu.Unlock()
	return nil
}

func (am *AuthenticationService) getSecret(role RoleName) (string, error) {
	envName, ok := ENV_NAME_BY_ROLE_NAME[role]
	if !ok {
		return "", fmt.Errorf("could not find env name for role %s", role)
	}
	secret, secretExists := os.LookupEnv(envName)
	if !secretExists {
		return "", fmt.Errorf("could not find bot client secret")
	}
	return secret, nil
}
