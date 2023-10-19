package server

import (
	"net/http"
)
import "fmt"
import "github.com/gorilla/websocket"

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

func listenOnWsCon(con *websocket.Conn) {
	msg := AuthMessageContent{
		Topic:      AUTH,
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		con, _ := upgradeToWSCon(w, r)
		go listenOnWsCon(con)
	})

	err := http.ListenAndServe(":8080", nil)
	fmt.Println("asdf")
	if err != nil {
		fmt.Println(err)
		return
	}
}
