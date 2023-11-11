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
	inChannel  chan *Message
	outChannel chan *Message
	conn       *websocket.Conn
	cleanup    func(*UserClient)
}

func NewUserClient(conn *websocket.Conn, cleanup func(*UserClient)) *UserClient {
	pubKey, priKey := GenerateKeyset()
	inChannel := make(chan *Message)
	outChannel := make(chan *Message)

	uc := UserClient{
		active:     true,
		publicKey:  pubKey,
		privateKey: priKey,
		inChannel:  inChannel,
		outChannel: outChannel,
		conn:       conn,
		cleanup:    cleanup,
	}
	go uc.listenOnServerChannel()
	go uc.listenOnWebsocket()

	logManagerConfigBuilder := NewLogManagerConfigBuilder()
	logManagerConfigBuilder.WithDecorator(pubKey, ClientKeyLogDecorator)
	if GetLogManager().Config.IsEnvMuted(pubKey) {
		logManagerConfigBuilder.WithMutedEnv(pubKey)
	}
	logConfig := logManagerConfigBuilder.Build()
	GetLogManager().InjectConfig(logConfig)

	msg := &Message{
		Topic:       "auth",
		ContentType: CONTENT_TYPE_AUTH,
		Content: &AuthMessageContent{
			PublicKey:  pubKey,
			PrivateKey: priKey,
		},
	}
	sendAuthErr := uc.SendMessage(msg)
	if sendAuthErr != nil {
		GetLogManager().LogRed("client", fmt.Sprintf("error sending auth message to client: %s", sendAuthErr), ALL_BUT_TEST_ENV)
	}
	return &uc
}

func (uc *UserClient) PublicKey() string {
	return uc.publicKey
}

func (uc *UserClient) InChannel() chan *Message {
	return uc.inChannel
}

func (uc *UserClient) OutChannel() chan *Message {
	return uc.outChannel
}

func (uc *UserClient) listenOnServerChannel() {
	for {
		time.Sleep(time.Millisecond * 1)
		select {
		case message := <-uc.inChannel:
			sendErr := uc.SendMessage(message)
			if sendErr != nil {
				GetLogManager().LogRed("client", fmt.Sprintf("error sending message to client: %s", sendErr), ALL_BUT_TEST_ENV)
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
		GetLogManager().Log("client", "cannot listen on websocket, connection is nil", ALL_BUT_TEST_ENV)
		return
	}
	for {
		_, rawMsg, readErr := uc.conn.ReadMessage()
		if !uc.active {
			return
		}
		if readErr != nil {
			GetLogManager().LogRed("client", fmt.Sprintf("error reading message from websocket: %s", readErr))
			// assume all readErrs are disconnects
			_ = userClientsManager.RemoveClient(uc)
			return
		}
		GetLogManager().Log(uc.publicKey, ">> ", string(rawMsg))

		msg, unmarshalErr := UnmarshalToMessage(rawMsg)
		if unmarshalErr != nil {
			GetLogManager().Log("client", fmt.Sprintf("could not unmarshal message: %s", unmarshalErr))
			continue
		}
		authErr := GetAuthManager().ValidateAuthInMessage(msg)
		if authErr != nil {
			GetLogManager().Log("client", fmt.Sprintf("auth error: %s", authErr))
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
	GetLogManager().Log(uc.publicKey, "<< ", string(msgJson))
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
