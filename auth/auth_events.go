package auth

import (
	models "github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/CameronHonis/service"
)

const (
	ROLE_SWITCHED EventVariant = "ROLE_SWITCHED"
	CREDS_VETTED  EventVariant = "CREDS_VETTED"
	CREDS_CHANGED EventVariant = "CREDS_CHANGED"
	CREDS_REMOVED EventVariant = "CREDS_REMOVED"
)

type RoleSwitchedPayload struct {
	ClientKey models.Key
	Role      models.RoleName
}

type RoleSwitchedEvent struct{ Event }

func NewRoleSwitchedEvent(clientKey models.Key, role models.RoleName) *RoleSwitchedEvent {
	return &RoleSwitchedEvent{
		Event: *NewEvent(ROLE_SWITCHED, &RoleSwitchedPayload{
			ClientKey: clientKey,
			Role:      role,
		}),
	}
}

type CredsVettedPayload struct {
	ClientKey models.Key
	PriKey    models.Key
}

type CredsVettedEvent struct{ Event }

func NewCredsVettedEvent(clientKey models.Key, priKey models.Key) *CredsVettedEvent {
	return &CredsVettedEvent{
		Event: *NewEvent(CREDS_VETTED, &CredsVettedPayload{
			ClientKey: clientKey,
			PriKey:    priKey,
		}),
	}
}

type CredsChangedPayload struct {
	OldCreds *models.AuthCreds
	NewCreds *models.AuthCreds
}

type CredsChangedEvent struct{ Event }

func NewCredsChangedEvent(oldCreds *models.AuthCreds, newCreds *models.AuthCreds) *CredsChangedEvent {
	return &CredsChangedEvent{
		Event: *NewEvent(CREDS_CHANGED, &CredsChangedPayload{
			OldCreds: oldCreds,
			NewCreds: newCreds,
		}),
	}
}

type CredsRemovedPayload struct {
	ClientKey models.Key
}

type CredsRemovedEvent struct{ Event }

func NewCredsRemovedEvent(clientKey models.Key) *CredsRemovedEvent {
	return &CredsRemovedEvent{
		Event: *NewEvent(CREDS_REMOVED, &CredsRemovedPayload{
			ClientKey: clientKey,
		}),
	}
}
