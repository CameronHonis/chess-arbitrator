package server

import (
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

var userClientsManager *UserClientsManager

func listenOnWsCon(con *websocket.Conn) {
	msg := AuthMessageContent{
		PublicKey:  "asdf",
		PrivateKey: "some-private-key",
	}
	msgJsonBytes, jsonParseErr := msg.toJsonBytes()
	if jsonParseErr != nil {
		fmt.Println(jsonParseErr)
		return
	}
	writeBytesErr := con.WriteMessage(websocket.TextMessage, msgJsonBytes)
	if writeBytesErr != nil {
		fmt.Println(writeBytesErr)
		return
	}
	for {
		messageType, p, err := con.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(messageType, string(p))
		if err := con.WriteMessage(websocket.TextMessage, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func StartWSServer() {
	userClientsManager, _ = NewUserClientsManager()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, connErr := upgradeToWSCon(w, r)
		if connErr != nil {
			fmt.Println(connErr)
			return
		}
		_, addClientErr := userClientsManager.AddNewClient(conn)
		if addClientErr != nil {
			fmt.Println(addClientErr)
			return
		}
	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		fmt.Println("asdf")
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
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

func handlePrompt(prompt *Prompt) {

}
