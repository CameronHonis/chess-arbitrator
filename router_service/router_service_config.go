package router_service

import (
	"github.com/CameronHonis/service"
	"os"
	"strconv"
)

type RouterServiceConfig struct {
	Port uint
}

func NewRouterServiceConfig() *RouterServiceConfig {
	portEnvVal, _ := os.LookupEnv("LISTEN_PORT")
	var port uint
	if num, err := strconv.Atoi(portEnvVal); err == nil {
		port = uint(num)
	} else {
		port = 8080
	}
	return &RouterServiceConfig{
		Port: port,
	}
}

func (rc *RouterServiceConfig) MergeWith(other service.ConfigI) service.ConfigI {
	newConfig := *(other.(*RouterServiceConfig))
	return &newConfig
}
