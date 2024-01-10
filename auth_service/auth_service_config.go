package auth_service

import (
	. "github.com/CameronHonis/service"
)

type AuthServiceConfig struct {
	ConfigI
}

func NewAuthServiceConfig() *AuthServiceConfig {
	return &AuthServiceConfig{}
}
