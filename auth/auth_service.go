package auth

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/secrets_manager"
	"github.com/CameronHonis/log"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"github.com/CameronHonis/set"
	"strconv"
	"sync"
	"time"
)

type AuthenticationServiceI interface {
	service.ServiceI
	GetRole(clientKey models.Key) (models.RoleName, error)
	ClientKeysByRole(roleName models.RoleName) *set.Set[models.Key]
	BotClientExists() bool

	CreateNewClient() *models.AuthCreds
	SwitchRole(clientKey models.Key, roleName models.RoleName, secret string) error
	RemoveClient(clientKey models.Key)
	RefreshPrivateKey(clientKey models.Key, priKey models.Key) error

	VetAuthInMessage(msg *models.Message) error
	VetPrivateKey(pubKey models.Key, priKey models.Key) error
	VetClientForTopic(clientKey models.Key, topic models.MessageTopic) error

	StripAuthFromMessage(msg *models.Message)
}

type AuthenticationService struct {
	service.Service

	__dependencies__ marker.Marker
	LoggerService    log.LoggerServiceI
	SecretsManager   secrets_manager.SecretsManagerI

	__state__         marker.Marker
	authCredsByClient map[models.Key]*models.AuthCreds
	clientKeysByRole  map[models.RoleName]*set.Set[models.Key]
	mu                sync.Mutex
}

func NewAuthenticationService(config *AuthServiceConfig) *AuthenticationService {
	authService := &AuthenticationService{
		authCredsByClient: make(map[models.Key]*models.AuthCreds),
		clientKeysByRole:  make(map[models.RoleName]*set.Set[models.Key]),
	}
	authService.Service = *service.NewService(authService, config)
	return authService
}

func (am *AuthenticationService) GetRole(clientKey models.Key) (models.RoleName, error) {
	creds, credsErr := am.getCreds(clientKey)
	if credsErr != nil {
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

func (am *AuthenticationService) CreateNewClient() *models.AuthCreds {
	clientKey, priKey := GenerateKeyset()
	creds := models.NewAuthCreds(clientKey, priKey, models.PLEB)
	am.setCreds(creds)
	return creds
}

func (am *AuthenticationService) SwitchRole(clientKey models.Key, roleName models.RoleName, secret string) error {
	creds, credsErr := am.getCreds(clientKey)
	if credsErr != nil {
		return fmt.Errorf("could not switch role: %s", credsErr.Error())
	}
	// filter out unauthorized role switches
	switch roleName {
	case models.PLEB:
		break
	case models.BOT:
		if am.BotClientExists() {
			return fmt.Errorf("bot already exists")
		}
		if am.SecretsManager.ValidateSecret(models.SECRET_BOT_CLIENT_SECRET, secret) != nil {
			return fmt.Errorf("invalid secret")
		}
	}

	// assumed that role switch is permitted after this point
	newCreds := builders.NewAuthCredsBuilder().FromAuthCreds(*creds).WithRole(roleName).Build()
	am.setCreds(newCreds)

	go am.Dispatch(NewRoleSwitchedEvent(clientKey, roleName))
	return nil
}

func (am *AuthenticationService) RemoveClient(clientKey models.Key) {
	am.removeCreds(clientKey)
}

func (am *AuthenticationService) RefreshPrivateKey(clientKey models.Key, priKey models.Key) error {
	priKeyStaleAfterMinsStr, configErr := am.SecretsManager.GetSecret(models.SECRET_AUTH_KEY_MINS_TO_STALE)
	if configErr != nil {
		return fmt.Errorf("couldnt refresh private key: %s", configErr)
	}
	priKeyStaleAfterMin, floatParseErr := strconv.ParseFloat(priKeyStaleAfterMinsStr, 64)
	if floatParseErr != nil {
		return fmt.Errorf("couldnt refresh private key: %s", floatParseErr)
	}
	creds, credsErr := am.getCreds(clientKey)
	if credsErr != nil {
		return fmt.Errorf("couldnt refresh private key: %s", credsErr)
	}
	if creds.PrivateKey != priKey {
		return fmt.Errorf("couldnt refresh private key: invalid private key")
	}

	minsSinceIssued := time.Now().Sub(creds.PriKeyCreatedAt).Minutes()
	newPriKey := creds.PrivateKey
	if minsSinceIssued >= priKeyStaleAfterMin {
		newPriKey = GeneratePriKey()
	}
	newCreds := builders.NewAuthCredsBuilder().FromAuthCreds(*creds).WithPrivateKey(newPriKey).Build()
	am.setCreds(newCreds)

	return nil
}

func (am *AuthenticationService) VetAuthInMessage(msg *models.Message) error {
	isValidAuth := ValidatePrivateKey(msg.SenderKey, msg.PrivateKey)
	if !isValidAuth {
		return fmt.Errorf("invalid auth")
	}
	return nil
}

func (am *AuthenticationService) VetPrivateKey(pubKey models.Key, priKey models.Key) error {
	creds, credsErr := am.getCreds(pubKey)
	if credsErr != nil {
		return fmt.Errorf("client with key %s does not exist", pubKey)
	}
	if priKey != creds.PrivateKey {
		return fmt.Errorf("private keys do not match")
	}
	return nil
}

func (am *AuthenticationService) StripAuthFromMessage(msg *models.Message) {
	msg.PrivateKey = ""
}

func (am *AuthenticationService) VetClientForTopic(clientKey models.Key, topic models.MessageTopic) error {
	return nil
}

func (am *AuthenticationService) getCreds(clientKey models.Key) (*models.AuthCreds, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	creds, ok := am.authCredsByClient[clientKey]
	if !ok {
		return nil, fmt.Errorf("creds do not exist")
	}
	return creds, nil
}

func (am *AuthenticationService) setCreds(creds *models.AuthCreds) {
	am.mu.Lock()
	defer am.mu.Unlock()
	prevCreds, _ := am.authCredsByClient[creds.ClientKey]
	if prevCreds == nil || creds.Role != prevCreds.Role {
		if prevCreds != nil {
			if oldRoleClients, ok := am.clientKeysByRole[prevCreds.Role]; ok {
				oldRoleClients.Remove(creds.ClientKey)
			}
		}
		newRoleClientKeys, ok := am.clientKeysByRole[creds.Role]
		if !ok {
			newRoleClientKeys = set.EmptySet[models.Key]()
			am.clientKeysByRole[creds.Role] = newRoleClientKeys
		}
		newRoleClientKeys.Add(creds.ClientKey)

		if prevCreds != nil {
			go am.Dispatch(NewRoleSwitchedEvent(creds.ClientKey, creds.Role))
		}
	}

	am.authCredsByClient[creds.ClientKey] = creds
	go am.Dispatch(NewCredsChangedEvent(prevCreds, creds))
}

func (am *AuthenticationService) removeCreds(clientKey models.Key) *models.AuthCreds {
	am.mu.Lock()
	defer am.mu.Unlock()
	creds, ok := am.authCredsByClient[clientKey]
	if !ok {
		return nil
	}
	delete(am.authCredsByClient, clientKey)
	go am.Dispatch(NewCredsRemovedEvent(clientKey))
	return creds
}
