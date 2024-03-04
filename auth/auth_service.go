package auth

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/secrets_manager"
	"github.com/CameronHonis/log"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"github.com/CameronHonis/set"
	"sync"
)

type AuthenticationServiceI interface {
	service.ServiceI
	GetRole(clientKey models.Key) (models.RoleName, error)
	ClientKeysByRole(roleName models.RoleName) *set.Set[models.Key]
	BotClientExists() bool

	CreateNewClient() *models.ClientAuthCreds
	SwitchRole(clientKey models.Key, roleName models.RoleName, secret string) error
	RemoveClient(clientKey models.Key) error
	RefreshPrivateKey(clientKey models.Key) (models.Key, error)

	ValidateAuthInMessage(msg *models.Message) error
	ValidatePrivateKey(pubKey models.Key, priKey models.Key) error
	StripAuthFromMessage(msg *models.Message)
	ValidateClientForTopic(clientKey models.Key, topic models.MessageTopic) error
}

type AuthenticationService struct {
	service.Service

	__dependencies__ marker.Marker
	LoggerService    log.LoggerServiceI
	SecretsManager   secrets_manager.SecretsManagerI

	__state__         marker.Marker
	authCredsByClient map[models.Key]*models.ClientAuthCreds
	clientKeysByRole  map[models.RoleName]*set.Set[models.Key]
	mu                sync.Mutex
}

func NewAuthenticationService(config *AuthServiceConfig) *AuthenticationService {
	authService := &AuthenticationService{
		authCredsByClient: make(map[models.Key]*models.ClientAuthCreds),
		clientKeysByRole:  make(map[models.RoleName]*set.Set[models.Key]),
	}
	authService.Service = *service.NewService(authService, config)
	return authService
}

func (am *AuthenticationService) GetRole(clientKey models.Key) (models.RoleName, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	creds, ok := am.authCredsByClient[clientKey]
	if !ok {
		return "", fmt.Errorf("could not find role for client %s", clientKey)
	}
	return creds.Role, nil
}

func (am *AuthenticationService) ClientKeysByRole(roleName models.RoleName) *set.Set[models.Key] {
	am.mu.Lock()
	clientKeys, ok := am.clientKeysByRole[roleName]
	am.mu.Unlock()
	if !ok {
		return set.EmptySet[models.Key]()
	}
	return clientKeys
}

func (am *AuthenticationService) BotClientExists() bool {
	botClientKeys := am.ClientKeysByRole(models.BOT)
	return botClientKeys.Size() > 0
}

func (am *AuthenticationService) CreateNewClient() *models.ClientAuthCreds {
	am.mu.Lock()
	defer am.mu.Unlock()

	clientKey, priKey := GenerateKeyset()
	creds := models.NewClientAuthCreds(clientKey, priKey, models.PLEB)
	am.authCredsByClient[clientKey] = creds
	clientKeys, ok := am.clientKeysByRole[models.PLEB]
	if !ok {
		clientKeys = set.EmptySet[models.Key]()
		am.clientKeysByRole[models.PLEB] = clientKeys
	}
	clientKeys.Add(clientKey)

	return creds
}

func (am *AuthenticationService) SwitchRole(clientKey models.Key, roleName models.RoleName, secret string) error {
	// filter out unauthorized role switches
	switch roleName {
	case models.PLEB:
		break
	case models.BOT:
		if am.BotClientExists() {
			return fmt.Errorf("bot already exists")
		}
		if am.SecretsManager.ValidateSecret(models.BOT_CLIENT_SECRET, secret) != nil {
			return fmt.Errorf("invalid secret")
		}
	}

	// assumed that role switch is permitted after this point
	roleErr := am.SetRole(clientKey, roleName)
	if roleErr != nil {
		go am.Dispatch(NewRoleSwitchDeniedEvent(clientKey, roleErr.Error()))
		return roleErr
	}

	go am.Dispatch(NewRoleSwitchGrantedEvent(clientKey, roleName))
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

func (am *AuthenticationService) RefreshPrivateKey(clientKey models.Key) (models.Key, error) {
}

func (am *AuthenticationService) ValidateAuthInMessage(msg *models.Message) error {
	isValidAuth := ValidatePrivateKey(msg.SenderKey, msg.PrivateKey)
	if !isValidAuth {
		return fmt.Errorf("invalid auth")
	}
	return nil
}

func (am *AuthenticationService) ValidatePrivateKey(pubKey models.Key, priKey models.Key) error {
}

func (am *AuthenticationService) StripAuthFromMessage(msg *models.Message) {
	msg.PrivateKey = ""
}

func (am *AuthenticationService) ValidateClientForTopic(clientKey models.Key, topic models.MessageTopic) error {
	return nil
}

func (am *AuthenticationService) SetRole(clientKey models.Key, role models.RoleName) error {
	oldRole, getRoleErr := am.GetRole(clientKey)
	if getRoleErr != nil {
		return getRoleErr
	}
	am.mu.Lock()
	am.roleByClient[clientKey] = role
	if oldRoleClients, ok := am.clientKeysByRole[oldRole]; ok {
		oldRoleClients.Remove(clientKey)
	}
	newRoleClientKeys, ok := am.clientKeysByRole[role]
	if !ok {
		newRoleClientKeys = set.EmptySet[models.Key]()
		am.clientKeysByRole[role] = newRoleClientKeys
	}
	newRoleClientKeys.Add(clientKey)
	am.mu.Unlock()
	return nil
}
