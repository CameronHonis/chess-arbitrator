package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type UserClient struct {
	publicKey     string
	privateKey    string
	serverChannel chan *Prompt
	stdoutMutex   *sync.Mutex
	conn          *websocket.Conn
	cleanup       func(*UserClient)
}

func NewUserClient(stdoutMutex *sync.Mutex, ch chan *Prompt, conn *websocket.Conn, cleanup func(*UserClient)) *UserClient {
	pubKey, priKey := generateKeyset()
	uc := UserClient{
		publicKey:     pubKey,
		privateKey:    priKey,
		serverChannel: ch,
		stdoutMutex:   stdoutMutex,
		conn:          conn,
		cleanup:       cleanup,
	}
	clientInitPrompt := Prompt{
		Type: PROMPT_TYPE_INIT_CLIENT,
		Content: &InitClientPromptContent{
			PublicKey: pubKey,
		},
	}
	ch <- &clientInitPrompt
	uc.listenOnServerChannel()
	uc.listenOnWebsocket()
	return &uc
}

func (uc *UserClient) listenOnServerChannel() {
	for {
		msg := <-uc.serverChannel
		//msgJson, jsonParseErr := json.Marshal(msg)
		//if jsonParseErr != nil {
		//	log.Fatal(jsonParseErr)
		//}
		//uc.stdoutMutex.Lock()
		//fmt.Println("Arbitrator >> ", string(msgJson))
		//uc.stdoutMutex.Unlock()
		uc.stdoutMutex.Lock()
		fmt.Println("Server >> ", msg)
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
		uc.stdoutMutex.Lock()
		fmt.Println(uc.publicKey, " >> ", string(msg))
		uc.stdoutMutex.Unlock()
	}
}
