package server

import (
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

type RouterConfig struct {
	Port uint
}

func NewRouterConfig() *RouterConfig {
	return &RouterConfig{
		Port: 8080,
	}
}

func (rc *RouterConfig) MergeWith(other ConfigI) ConfigI {
	newConfig := *(other.(*RouterConfig))
	return &newConfig
}

type RouterService struct {
	Service[*RouterConfig]

	__dependencies__   Marker
	UserClientsService UserClientsServiceI

	__state__ Marker
}

func NewRouterService(config *RouterConfig) *RouterService {
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

	config := rs.Config().(*RouterConfig)
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
