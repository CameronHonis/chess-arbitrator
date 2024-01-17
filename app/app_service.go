package app

import (
	"github.com/CameronHonis/chess-arbitrator/router_service"
	. "github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
)

//go:generate mockgen -destination mock/app_service_mock.go . AppServiceI
type AppServiceI interface {
	service.ServiceI
}

type AppService struct {
	service.Service

	__dependencies__ Marker
	RouterService    router_service.RouterServiceI

	__state__ Marker
}

func NewAppService(config *AppServiceConfig) *AppService {
	app := &AppService{}
	app.Service = *service.NewService(app, config)
	return app
}
