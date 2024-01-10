package router_service

import (
	"github.com/CameronHonis/chess-arbitrator/user_clients_service"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

type RouterServiceConfig struct {
	Port uint
}

func NewRouterServiceConfig() *RouterServiceConfig {
	return &RouterServiceConfig{
		Port: 8080,
	}
}

func (rc *RouterServiceConfig) MergeWith(other ConfigI) ConfigI {
	newConfig := *(other.(*RouterServiceConfig))
	return &newConfig
}

type RouterServiceI interface {
	ServiceI
	StartWSServer()
}

type RouterService struct {
	Service

	__dependencies__   Marker
	UserClientsService user_clients_service.UserClientsServiceI

	__state__ Marker
}

func NewRouterService(config *RouterServiceConfig) *RouterService {
	routerService := &RouterService{}
	routerService.Service = *NewService(routerService, config)
	return routerService
}

func (rs *RouterService) StartWSServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, connErr := upgradeToWSCon(w, r)
		if connErr != nil {
			fmt.Println(connErr)
			return
		}
		_, addClientErr := rs.UserClientsService.AddNewClient(conn)
		if addClientErr != nil {
			fmt.Println(addClientErr)
			return
		}
	})

	config := rs.Config().(*RouterServiceConfig)
	port := config.Port
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
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
