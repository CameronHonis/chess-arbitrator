package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"os"
	"sync"
)

const (
	AUTH_GRANTED EventVariant = "AUTH_GRANTED"
	AUTH_DENIED               = "AUTH_DENIED"
)

type AuthenticationGrantedPayload struct {
	ClientKey Key
	Role      RoleName
}

type AuthenticationGrantedEvent struct{ Event }

func NewAuthenticationGrantedEvent(clientKey Key, role RoleName) *AuthenticationGrantedEvent {
	return &AuthenticationGrantedEvent{
		Event: *NewEvent(AUTH_GRANTED, &AuthenticationGrantedPayload{
			ClientKey: clientKey,
			Role:      role,
		}),
	}
}

type AuthenticationDeniedPayload struct {
	ClientKey Key
	Reason    string
}

type AuthenticationDeniedEvent struct{ Event }

func NewAuthenticationDeniedEvent(clientKey Key, reason string) *AuthenticationDeniedEvent {
	return &AuthenticationDeniedEvent{
		Event: *NewEvent(AUTH_DENIED, &AuthenticationDeniedPayload{
			ClientKey: clientKey,
			Reason:    reason,
		}),
	}
}

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

func NewAuthenticationConfig() *AuthConfig {
	return &AuthConfig{}
}

type AuthenticationServiceI interface {
	ServiceI
	GetRole(clientKey Key) (RoleName, error)

	AddClient(clientKey Key)
	UpgradeAuth(clientKey Key, roleName RoleName, secret string) error
	RemoveClient(clientKey Key) error

	ValidateSecret(roleName RoleName, secret string) error
	ValidateAuthInMessage(msg *Message) error
	StripAuthFromMessage(msg *Message)
	ValidateClientForTopic(clientKey Key, topic MessageTopic) error
}

type AuthenticationService struct {
	Service

	__dependencies__ Marker
	LoggerService    LoggerServiceI

	__state__    Marker
	roleByClient map[Key]RoleName
	mu           sync.Mutex
}

func NewAuthenticationService(config *AuthConfig) *AuthenticationService {
	authService := &AuthenticationService{}
	authService.Service = *NewService(authService, config)
	return authService
}

func (am *AuthenticationService) GetRole(clientKey Key) (RoleName, error) {
	am.mu.Lock()
	role, ok := am.roleByClient[clientKey]
	am.mu.Unlock()
	if !ok {
		return "", fmt.Errorf("could not find role for client %s", clientKey)
	}
	return role, nil
}

func (am *AuthenticationService) AddClient(clientKey Key) {
	am.mu.Lock()
	am.roleByClient[clientKey] = PLEB
	am.mu.Unlock()
}

func (am *AuthenticationService) UpgradeAuth(clientKey Key, roleName RoleName, secret string) error {
	validSecretErr := am.ValidateSecret(roleName, secret)
	if validSecretErr != nil {
		am.Dispatch(NewAuthenticationDeniedEvent(clientKey, validSecretErr.Error()))
		return validSecretErr
	}

	roleErr := am.SetRole(clientKey, roleName)
	if roleErr != nil {
		am.Dispatch(NewAuthenticationDeniedEvent(clientKey, roleErr.Error()))
		return roleErr
	}

	am.Dispatch(NewAuthenticationGrantedEvent(clientKey, roleName))
	return nil
}

func (am *AuthenticationService) RemoveClient(clientKey Key) error {
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
	expSecret, secretErr := am.GetSecret(roleName)
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

func (am *AuthenticationService) StripAuthFromMessage(msg *Message) {
	msg.PrivateKey = ""
}

func (am *AuthenticationService) ValidateClientForTopic(clientKey Key, topic MessageTopic) error {
	return nil
}

func (am *AuthenticationService) SetRole(clientKey Key, role RoleName) error {
	_, getRoleErr := am.GetRole(clientKey)
	if getRoleErr != nil {
		return getRoleErr
	}
	am.mu.Lock()
	am.roleByClient[clientKey] = role
	am.mu.Unlock()
	return nil
}

func (am *AuthenticationService) GetSecret(role RoleName) (string, error) {
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
