package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type UserClient struct {
	active        bool //assumed that cleanup already ran if set to true
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
		active:        true,
		publicKey:     pubKey,
		privateKey:    priKey,
		serverChannel: ch,
		stdoutMutex:   stdoutMutex,
		conn:          conn,
		cleanup:       cleanup,
	}
	//clientInitPrompt := Prompt{
	//	Type: PROMPT_TYPE_INIT_CLIENT,
	//	Content: &InitClientPromptContent{
	//		PublicKey: pubKey,
	//	},
	//}
	//ch <- &clientInitPrompt
	go uc.listenOnServerChannel()
	go uc.listenOnWebsocket()
	return &uc
}

func (uc *UserClient) PublicKey() string {
	return uc.publicKey
}

func (uc *UserClient) ServerChannel() chan *Prompt {
	return uc.ServerChannel()
}

func (uc *UserClient) listenOnServerChannel() {
	for {
		select {
		case prompt := <-uc.serverChannel:
			uc.handlePrompt(prompt)
		default:
			if !uc.active {
				return
			}
		}
	}
}

func (uc *UserClient) listenOnWebsocket() {
	if uc.conn == nil {
		fmt.Println("cannot listen on websocket, connection is nil")
		return
	}
	for {
		_, rawMsg, readErr := uc.conn.ReadMessage()
		if !uc.active {
			return
		}
		if readErr != nil {
			// Possibly handle non-connection closing read errors?
			// For now, assume all errors are connection closing
			uc.Kill()
			return
		}
		uc.stdoutMutex.Lock()
		fmt.Println(uc.publicKey, " >> ", string(rawMsg))
		uc.stdoutMutex.Unlock()

		msg, unmarshalErr := UnmarshalToMessage(rawMsg)
		if unmarshalErr != nil {
			uc.stdoutMutex.Lock()
			fmt.Println(unmarshalErr)
			uc.stdoutMutex.Unlock()
			continue
		}
		uc.handleMessage(msg)
	}
}

func (uc *UserClient) handlePrompt(prompt *Prompt) {
	switch prompt.Type {
	case PROMPT_TYPE_INIT_CLIENT:

	}
}

func (uc *UserClient) handleMessage(msg *Message) {

}

func (uc *UserClient) Kill() {
	uc.cleanup(uc)

	uc.active = false
}
