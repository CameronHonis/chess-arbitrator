package clients_manager

import "github.com/CameronHonis/service"

type ClientsManagerConfig struct {
	service.ConfigI
}

func NewClientsManagerConfig() *ClientsManagerConfig {
	return &ClientsManagerConfig{}
}
