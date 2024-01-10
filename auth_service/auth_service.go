package auth_service

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/helpers"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"os"
	"sync"
)

type AuthenticationServiceI interface {
	ServiceI
	GetRole(clientKey models.Key) (models.RoleName, error)

	AddClient(clientKey models.Key)
	UpgradeAuth(clientKey models.Key, roleName models.RoleName, secret string) error
	RemoveClient(clientKey models.Key) error

	ValidateSecret(roleName models.RoleName, secret string) error
	ValidateAuthInMessage(msg *models.Message) error
	StripAuthFromMessage(msg *models.Message)
	ValidateClientForTopic(clientKey models.Key, topic models.MessageTopic) error
}

type AuthenticationService struct {
	Service

	__dependencies__ Marker
	LoggerService    LoggerServiceI

	__state__    Marker
	roleByClient map[models.Key]models.RoleName
	mu           sync.Mutex
}

func NewAuthenticationService(config *AuthServiceConfig) *AuthenticationService {
	authService := &AuthenticationService{}
	authService.Service = *NewService(authService, config)
	return authService
}

func (am *AuthenticationService) GetRole(clientKey models.Key) (models.RoleName, error) {
	am.mu.Lock()
	role, ok := am.roleByClient[clientKey]
	am.mu.Unlock()
	if !ok {
		return "", fmt.Errorf("could not find role for client %s", clientKey)
	}
	return role, nil
}

func (am *AuthenticationService) AddClient(clientKey models.Key) {
	am.mu.Lock()
	am.roleByClient[clientKey] = models.PLEB
	am.mu.Unlock()
}

func (am *AuthenticationService) UpgradeAuth(clientKey models.Key, roleName models.RoleName, secret string) error {
	validSecretErr := am.ValidateSecret(roleName, secret)
	if validSecretErr != nil {
		go am.Dispatch(NewAuthenticationDeniedEvent(clientKey, validSecretErr.Error()))
		return validSecretErr
	}

	roleErr := am.SetRole(clientKey, roleName)
	if roleErr != nil {
		go am.Dispatch(NewAuthenticationDeniedEvent(clientKey, roleErr.Error()))
		return roleErr
	}

	go am.Dispatch(NewAuthenticationGrantedEvent(clientKey, roleName))
	return nil
}

func (am *AuthenticationService) RemoveClient(clientKey models.Key) error {
	_, roleErr := am.GetRole(clientKey)
	if roleErr != nil {
		am.LoggerService.LogRed(models.ENV_SERVER, fmt.Sprintf("attempt to remove non-existant client %s", clientKey))
		return nil
	}
	am.mu.Lock()
	delete(am.roleByClient, clientKey)
	am.mu.Unlock()
	return nil
}

func (am *AuthenticationService) ValidateSecret(roleName models.RoleName, secret string) error {
	expSecret, secretErr := am.GetSecret(roleName)
	if secretErr != nil {
		return secretErr
	}
	if secret != expSecret {
		return fmt.Errorf("invalid secret")
	}
	return nil
}

func (am *AuthenticationService) ValidateAuthInMessage(msg *models.Message) error {
	isValidAuth := helpers.ValidatePrivateKey(msg.SenderKey, msg.PrivateKey)
	if !isValidAuth {
		return fmt.Errorf("invalid auth")
	}
	return nil
}

func (am *AuthenticationService) StripAuthFromMessage(msg *models.Message) {
	msg.PrivateKey = ""
}

func (am *AuthenticationService) ValidateClientForTopic(clientKey models.Key, topic models.MessageTopic) error {
	return nil
}

func (am *AuthenticationService) SetRole(clientKey models.Key, role models.RoleName) error {
	_, getRoleErr := am.GetRole(clientKey)
	if getRoleErr != nil {
		return getRoleErr
	}
	am.mu.Lock()
	am.roleByClient[clientKey] = role
	am.mu.Unlock()
	return nil
}

func (am *AuthenticationService) GetSecret(role models.RoleName) (string, error) {
	envName, ok := models.ENV_NAME_BY_ROLE_NAME[role]
	if !ok {
		return "", fmt.Errorf("could not find env name for role %s", role)
	}
	secret, secretExists := os.LookupEnv(envName)
	if !secretExists {
		return "", fmt.Errorf("could not find bot client secret")
	}
	return secret, nil
}
