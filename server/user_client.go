package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Message struct {
	Topic   MessageTopic
	Content interface{}
}

type UserClient struct {
	publicKey   string
	privateKey  string
	channel     chan Message
	stdoutMutex *sync.Mutex
	conn        *websocket.Conn
	cleanup     func(*UserClient)
}

func NewUserClient(stdoutMutex *sync.Mutex, ch chan Message, conn *websocket.Conn, cleanup func(*UserClient)) *UserClient {
	pubKey, priKey := generateKeyset()
	uc := UserClient{
		publicKey:   pubKey,
		privateKey:  priKey,
		channel:     ch,
		stdoutMutex: stdoutMutex,
		conn:        conn,
		cleanup:     cleanup,
	}
	uc.listenOnChannel()
	uc.listenOnWebsocket()
	return &uc
}

func (uc *UserClient) listenOnChannel() {
	for {
		msg := <-uc.channel
		msgJson, jsonParseErr := json.Marshal(msg)
		if jsonParseErr != nil {
			log.Fatal(jsonParseErr)
		}
		uc.stdoutMutex.Lock()
		fmt.Println("Arbitrator >> ", string(msgJson))
		uc.stdoutMutex.Unlock()
	}
}

func (uc *UserClient) listenOnWebsocket() {
	for {
		_, msg, readErr := uc.conn.ReadMessage()
		if readErr != nil {
			// Possibly handle non-connection closing read errors?
			// For now, assume all errors are connection closing
			uc.cleanup(uc)
			return
		}

	}
}
