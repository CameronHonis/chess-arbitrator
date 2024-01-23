package clients_manager

import "github.com/CameronHonis/service"

type ClientsManagerConfig struct {
	service.ConfigI
	EventListenersAttached bool
}

func NewClientsManagerConfig(eventListenersAttached bool) *ClientsManagerConfig {
	return &ClientsManagerConfig{
		EventListenersAttached: eventListenersAttached,
	}
}
