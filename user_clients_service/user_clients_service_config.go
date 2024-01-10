package user_clients_service

import "github.com/CameronHonis/service"

type UserClientsServiceConfig struct {
	service.ConfigI
}

func NewUserClientsServiceConfig() *UserClientsServiceConfig {
	return &UserClientsServiceConfig{}
}
