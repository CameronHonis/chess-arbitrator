package auth

import (
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
)

const (
	AUTH_UPGRADE_GRANTED EventVariant = "AUTH_UPGRADE_GRANTED"
	AUTH_UPGRADE_DENIED               = "AUTH_UPGRADE_DENIED"
)

type AuthUpgradeGrantedPayload struct {
	ClientKey Key
	Role      RoleName
}

type AuthUpgradeGrantedEvent struct{ Event }

func NewAuthUpgradeGrantedEvent(clientKey Key, role RoleName) *AuthUpgradeGrantedEvent {
	return &AuthUpgradeGrantedEvent{
		Event: *NewEvent(AUTH_UPGRADE_GRANTED, &AuthUpgradeGrantedPayload{
			ClientKey: clientKey,
			Role:      role,
		}),
	}
}

type AuthUpgradeDeniedPayload struct {
	ClientKey Key
	Reason    string
}

type AuthUpgradeDeniedEvent struct{ Event }

func NewAuthUpgradeDeniedEvent(clientKey Key, reason string) *AuthUpgradeDeniedEvent {
	return &AuthUpgradeDeniedEvent{
		Event: *NewEvent(AUTH_UPGRADE_DENIED, &AuthUpgradeDeniedPayload{
			ClientKey: clientKey,
			Reason:    reason,
		}),
	}
}
