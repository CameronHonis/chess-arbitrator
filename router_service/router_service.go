package router_service

import (
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/user_clients_service"
	"github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

type RouterServiceI interface {
	ServiceI
	StartWSServer()
}

type RouterService struct {
	Service

	__dependencies__   Marker
	UserClientsService user_clients_service.UserClientsServiceI
	Logger             log.LoggerServiceI

	__state__ Marker
}

func NewRouterService(config *RouterServiceConfig) *RouterService {
	routerService := &RouterService{}
	routerService.Service = *NewService(routerService, config)
	return routerService
}

func (rs *RouterService) OnStart() {
	go rs.StartWSServer()
}

func (rs *RouterService) StartWSServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, connErr := upgradeToWSCon(w, r)
		if connErr != nil {
			rs.Logger.LogRed(models.ENV_SERVER, "error upgrading to ws conn:", connErr.Error())
			return
		}
		_, addClientErr := rs.UserClientsService.AddNewClient(conn)
		if addClientErr != nil {
			rs.Logger.LogRed(models.ENV_SERVER, "error adding client:", connErr.Error())
			return
		}
	})

	config := rs.Config().(*RouterServiceConfig)
	port := config.Port
	addr := fmt.Sprintf(":%d", port)
	rs.Logger.Log(models.ENV_SERVER, "server spinning up on port ", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		rs.Logger.LogRed(models.ENV_SERVER, "could not spin up server:", err.Error())
		return
	}
}

func upgradeToWSCon(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return con, nil
}
