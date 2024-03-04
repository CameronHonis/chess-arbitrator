package models

import "time"

type AuthCreds struct {
	ClientKey     Key
	PrivateKey    Key
	CreatedAt     time.Time
	PriKeyStaleAt time.Time
	Role          RoleName
}

func NewAuthCreds(clientKey Key, privateKey Key, role RoleName) *AuthCreds {
	return &AuthCreds{
		ClientKey:     clientKey,
		PrivateKey:    privateKey,
		CreatedAt:     time.Now(),
		PriKeyStaleAt: time.Now().Add(time.Hour * 24 * 7),
		Role:          role,
	}
}
