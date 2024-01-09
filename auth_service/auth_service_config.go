package auth_service

import (
	. "github.com/CameronHonis/service"
)

type AuthConfig struct {
	ConfigI
}

func NewAuthenticationConfig() *AuthConfig {
	return &AuthConfig{}
}
