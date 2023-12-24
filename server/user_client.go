package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

type UserClient struct {
	active     bool //assumed that cleanup already ran if set to true
	publicKey  string
	privateKey string
	conn       *websocket.Conn
	cleanup    func(*UserClient)
}

func NewUserClient(conn *websocket.Conn, cleanup func(*UserClient)) *UserClient {
	pubKey, priKey := GenerateKeyset()

	uc := &UserClient{
		active:     true,
		publicKey:  pubKey,
		privateKey: priKey,
		conn:       conn,
		cleanup:    cleanup,
	}
	return uc
}

func (uc *UserClient) Active() bool {
	return uc.active
}

func (uc *UserClient) PublicKey() string {
	return uc.publicKey
}

func (uc *UserClient) WSConn() *websocket.Conn {
	return uc.conn
}

func (uc *UserClient) listenOnServerChannel() {
	for {
		time.Sleep(time.Millisecond * 1)
		select {
		case message := <-uc.inChannel:
			sendErr := uc.SendMessage(message)
			if sendErr != nil {
				uc.logManager.LogRed(ENV_CLIENT, fmt.Sprintf("error sending message to client: %s", sendErr), ALL_BUT_TEST_ENV)
			}
		default:
			if !uc.active {
				return
			}
		}
	}
}

func (uc *UserClient) listenOnWebsocket() {
	if uc.conn == nil {
		uc.logManager.Log(ENV_CLIENT, "cannot listen on websocket, connection is nil", ALL_BUT_TEST_ENV)
		return
	}
	for {
		_, rawMsg, readErr := uc.conn.ReadMessage()
		if !uc.active {
			return
		}
		if readErr != nil {
			uc.logManager.LogRed(ENV_CLIENT, fmt.Sprintf("error reading message from websocket: %s", readErr))
			// assume all readErrs are disconnects
			_ = userClientsManager.RemoveClient(uc)
			return
		}
		uc.logManager.Log(uc.publicKey, ">> ", string(rawMsg))

		msg, unmarshalErr := UnmarshalToMessage(rawMsg)
		if unmarshalErr != nil {
			uc.logManager.Log(ENV_CLIENT, fmt.Sprintf("could not unmarshal message: %s", unmarshalErr))
			continue
		}
		authErr := uc.authManager.ValidateAuthInMessage(msg)
		if authErr != nil {
			uc.logManager.Log(ENV_CLIENT, fmt.Sprintf("auth error: %s", authErr))
			continue
		}
		uc.outChannel <- msg
	}
}

func (uc *UserClient) SendMessage(msg *Message) error {
	if uc.conn == nil {
		return fmt.Errorf("cannot send message, connection is nil: %s", msg)
	}
	msgJson, jsonErr := msg.Marshal()
	if jsonErr != nil {
		return jsonErr
	}
	uc.logManager.Log(uc.publicKey, "<< ", string(msgJson))
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

func ClientKeyLogDecorator(clientKey string) string {
	concatKey := clientKey[:4] + ".." + clientKey[len(clientKey)-4:]
	return WrapCyan(fmt.Sprintf("%s", strings.ToUpper(concatKey)))
}
