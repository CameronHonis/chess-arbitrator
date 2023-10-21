package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type UserClient struct {
	active      bool //assumed that cleanup already ran if set to true
	publicKey   string
	privateKey  string
	inChannel   chan *Prompt
	outChannel  chan *Prompt
	stdoutMutex *sync.Mutex
	conn        *websocket.Conn
	cleanup     func(*UserClient)
}

func NewUserClient(stdoutMutex *sync.Mutex, conn *websocket.Conn, cleanup func(*UserClient)) *UserClient {
	pubKey, priKey := generateKeyset()
	inChannel := make(chan *Prompt)
	outChannel := make(chan *Prompt)

	uc := UserClient{
		active:      true,
		publicKey:   pubKey,
		privateKey:  priKey,
		inChannel:   inChannel,
		outChannel:  outChannel,
		stdoutMutex: stdoutMutex,
		conn:        conn,
		cleanup:     cleanup,
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

func (uc *UserClient) InChannel() chan *Prompt {
	return uc.inChannel
}

func (uc *UserClient) OutChannel() chan *Prompt {
	return uc.outChannel
}

func (uc *UserClient) listenOnServerChannel() {
	for {
		time.Sleep(time.Millisecond * 1)
		select {
		case prompt := <-uc.inChannel:
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
			fmt.Println("error reading message from websocket: ", readErr)
			// assume all readErrs are disconnects
			userClientsManager.RemoveClient(uc)
			return
		}
		if LOG_INCOMING_MESSAGES {
			uc.stdoutMutex.Lock()
			fmt.Println("[CLIENT | ", uc.publicKey, "] >> ", string(rawMsg))
			uc.stdoutMutex.Unlock()
		}

		//msg, unmarshalErr := UnmarshalToMessage(rawMsg)
		//if unmarshalErr != nil {
		//	uc.stdoutMutex.Lock()
		//	fmt.Println(unmarshalErr)
		//	uc.stdoutMutex.Unlock()
		//	continue
		//}
		//uc.handleMessage(msg)
	}
}

func (uc *UserClient) SendMessage(msg *Message) error {
	msgJson, jsonErr := msg.Marshal()
	if jsonErr != nil {
		return jsonErr
	}
	writeErr := uc.conn.WriteMessage(websocket.TextMessage, msgJson)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func (uc *UserClient) Kill() {
	uc.cleanup(uc)

	uc.active = false
}

func (uc *UserClient) handlePrompt(prompt *Prompt) {
	switch prompt.Type {
	case PROMPT_TYPE_TRANSFER_MESSAGE:
	}
}

func (uc *UserClient) handleMessage(msg *Message) {
	switch msg.Topic {
	case MESSAGE_TOPIC_INIT_BOT_MATCH:
		uc.handleInitBotMatchMessage(msg.Content.(*InitBotMatchMessageContent))
	}
	if !msg.IsPrivate() {
		uc.outChannel <- &Prompt{
			Type:      PROMPT_TYPE_TRANSFER_MESSAGE,
			SenderKey: uc.publicKey,
			Content: &TransferMessagePromptContent{
				Message: msg,
			},
		}
	}
}

func (uc *UserClient) handleTransferMessagePrompt(content *TransferMessagePromptContent) {

}

func (uc *UserClient) handleInitBotMatchMessage(content *InitBotMatchMessageContent) {

}
