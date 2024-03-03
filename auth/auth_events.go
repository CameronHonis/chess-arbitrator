package auth

import (
	. "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
)

const (
	ROLE_SWITCH_GRANTED EventVariant = "ROLE_SWITCH_GRANTED"
	ROLE_SWITCH_DENIED  EventVariant = "ROLE_SWITCH_DENIED"
)

type RoleSwitchGrantedPayload struct {
	ClientKey Key
	Role      RoleName
}

type RoleSwitchGrantedEvent struct{ Event }

func NewRoleSwitchGrantedEvent(clientKey Key, role RoleName) *RoleSwitchGrantedEvent {
	return &RoleSwitchGrantedEvent{
		Event: *NewEvent(ROLE_SWITCH_GRANTED, &RoleSwitchGrantedPayload{
			ClientKey: clientKey,
			Role:      role,
		}),
	}
}

type RoleSwitchDeniedPayload struct {
	ClientKey Key
	Reason    string
}

type RoleSwitchDeniedEvent struct{ Event }

func NewRoleSwitchDeniedEvent(clientKey Key, reason string) *RoleSwitchDeniedEvent {
	return &RoleSwitchDeniedEvent{
		Event: *NewEvent(ROLE_SWITCH_DENIED, &RoleSwitchDeniedPayload{
			ClientKey: clientKey,
			Reason:    reason,
		}),
	}
}
