package models

import "time"

type AuthCreds struct {
	ClientKey       Key
	PrivateKey      Key
	CreatedAt       time.Time
	PriKeyCreatedAt time.Time
	Role            RoleName
}

func NewAuthCreds(clientKey Key, privateKey Key, role RoleName) *AuthCreds {
	return &AuthCreds{
		ClientKey:       clientKey,
		PrivateKey:      privateKey,
		CreatedAt:       time.Now(),
		PriKeyCreatedAt: time.Now(),
		Role:            role,
	}
}
