package sub_service

import (
	"github.com/CameronHonis/service"
)

type SubscriptionServiceConfig struct {
	service.ConfigI
}

func NewSubscriptionServiceConfig() *SubscriptionServiceConfig {
	return &SubscriptionServiceConfig{}
}
