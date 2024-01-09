package auth_service

import (
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
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
