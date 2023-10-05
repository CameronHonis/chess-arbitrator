package server

import "net/http"
import "fmt"
import "github.com/gorilla/websocket"

func listenOnWsCon(con *websocket.Conn) {
	for {
		messageType, p, err := con.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(messageType, string(p))
		if err := con.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func StartWSServer() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		con, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		listenOnWsCon(con)
	})

	http.ListenAndServe(":8080", nil)
}
