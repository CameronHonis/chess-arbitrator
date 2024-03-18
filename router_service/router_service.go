package router_service

import (
	"context"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/log"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

type RouterServiceI interface {
	service.ServiceI
	StartWSServer()
}

type RouterService struct {
	service.Service

	__dependencies__ marker.Marker
	ClientsManager   clients_manager.ClientsManagerI
	Logger           log.LoggerServiceI

	__state__ marker.Marker
	server    *http.Server
}

func NewRouterService(config *RouterServiceConfig) *RouterService {
	routerService := &RouterService{}
	routerService.Service = *service.NewService(routerService, config)
	return routerService
}

func (rs *RouterService) OnStart() {
	go rs.StartWSServer()
}

func (rs *RouterService) OnStop() {
	if err := rs.server.Shutdown(context.Background()); err != nil {
		rs.Logger.LogRed(models.ENV_SERVER, "could not stop server:", err)
	}
}

func (rs *RouterService) StartWSServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, connErr := upgradeToWSCon(w, r)
		if connErr != nil {
			rs.Logger.LogRed(models.ENV_SERVER, "error upgrading to ws conn:", connErr)
			return
		}
		rs.ClientsManager.AddConn(conn)
	})

	config := rs.Config().(*RouterServiceConfig)
	port := config.Port
	addr := fmt.Sprintf(":%d", port)

	rs.Logger.Log(models.ENV_SERVER, "server spinning up on port", port)
	rs.server = &http.Server{Addr: addr}
	err := rs.server.ListenAndServe()
	if err != nil {
		rs.Logger.LogRed(models.ENV_SERVER, "could not spin up server:", err)
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
