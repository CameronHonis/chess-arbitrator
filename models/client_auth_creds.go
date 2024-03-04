package models

import "time"

type ClientAuthCreds struct {
	ClientKey     Key
	PrivateKey    Key
	CreatedAt     time.Time
	PriKeyStaleAt time.Time
	Role          RoleName
}

func NewClientAuthCreds(clientKey Key, privateKey Key, role RoleName) *ClientAuthCreds {
	return &ClientAuthCreds{
		ClientKey:     clientKey,
		PrivateKey:    privateKey,
		CreatedAt:     time.Now(),
		PriKeyStaleAt: time.Now().Add(time.Hour * 24 * 7),
		Role:          role,
	}
}
