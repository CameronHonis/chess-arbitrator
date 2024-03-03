package secrets_manager

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
	"os"
)

type SecretsManagerI interface {
	service.ServiceI
	GetSecret(name models.Secret) (string, error)
	ValidateSecret(name models.Secret, secret string) error
}

type SecretsManager struct {
	service.Service
}

func NewSecretsManager() *SecretsManager {
	secretsManager := &SecretsManager{}
	secretsManager.Service = *service.NewService(secretsManager, nil)
	return secretsManager
}

func (sm *SecretsManager) GetSecret(name models.Secret) (string, error) {
	if secret, _ := os.LookupEnv(string(name)); secret != "" {
		return secret, nil
	}
	return "", fmt.Errorf("secret %s not configured", name)
}

func (sm *SecretsManager) ValidateSecret(name models.Secret, secret string) error {
	actualSecret, err := sm.GetSecret(name)
	if err != nil {
		return err
	}
	if actualSecret != secret {
		return fmt.Errorf("secret %s with name %s does not match", secret, name)
	}
	return nil
}
