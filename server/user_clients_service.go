package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"github.com/gorilla/websocket"
	"sync"
)

type UserClientsConfig struct {
}

func NewUserClientsConfig() *UserClientsConfig {
	return &UserClientsConfig{}
}

func (uc *UserClientsConfig) MergeWith(other ConfigI) ConfigI {
	newConfig := *(other.(*UserClientsConfig))
	return &newConfig
}

type UserClientsServiceI interface {
	ServiceI
	GetClient(clientKey string) (*UserClient, error)
	AddNewClient(conn *websocket.Conn) (*UserClient, error)
	AddClient(client *UserClient) error
	RemoveClient(client *UserClient) error
	BroadcastMessage(message *Message)
	DirectMessage(message *Message, clientKey string) error
}

type UserClientsService struct {
	Service[*UserClientsConfig]

	__dependencies__ Marker
	LoggerService    LoggerServiceI
	MsgService       MessageServiceI
	SubService       SubscriptionServiceI
	AuthService      AuthenticationServiceI

	__state__         Marker
	mu                sync.Mutex
	clientByPublicKey map[string]*UserClient
}

func NewUserClientsService(config *UserClientsConfig) *UserClientsService {
	userClientsService := &UserClientsService{
		clientByPublicKey: make(map[string]*UserClient),
	}
	userClientsService.Service = *NewService(userClientsService, config)
	return userClientsService
}

func (uc *UserClientsService) AddNewClient(conn *websocket.Conn) (*UserClient, error) {
	client := NewUserClient(conn, uc.CleanupClient)

	err := uc.AddClient(client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (uc *UserClientsService) AddClient(client *UserClient) error {
	uc.mu.Lock()
	if _, ok := uc.clientByPublicKey[client.PublicKey()]; ok {
		client.Kill()
		return fmt.Errorf("client with key %s already exists", client.PublicKey())
	}
	uc.clientByPublicKey[client.PublicKey()] = client
	uc.mu.Unlock()
	go uc.listenForUserInput(client)
	return nil
}

func (uc *UserClientsService) RemoveClient(client *UserClient) error {
	pubKey := client.PublicKey()

	uc.mu.Lock()
	if _, ok := uc.clientByPublicKey[pubKey]; !ok {
		uc.mu.Unlock()
		return fmt.Errorf("client with key %s isn't an established client", client.PublicKey())
	}
	delete(uc.clientByPublicKey, pubKey)
	uc.mu.Unlock()

	client.Kill()
	uc.SubService.UnsubClientFromAll(pubKey)
	return nil
}

func (uc *UserClientsService) GetClient(clientKey string) (*UserClient, error) {
	defer uc.mu.Unlock()
	uc.mu.Lock()
	client, ok := uc.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("no client with public key %s", clientKey)
	}
	return client, nil
}

func (uc *UserClientsService) listenForUserInput(userClient *UserClient) {
	for {
		_, rawMsg, readErr := userClient.WSConn().ReadMessage()
		if !userClient.Active() {
			return
		}
		if readErr != nil {
			uc.LoggerService.LogRed(ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", readErr), ALL_BUT_TEST_ENV)
			// assume all readErrs are disconnects
			_ = uc.RemoveClient(userClient)
			return
		}
		if err := uc.readMessage(userClient.PublicKey(), rawMsg); err != nil {
			uc.LoggerService.LogRed(ENV_SERVER, fmt.Sprintf("error reading message from websocket: %s", err), ALL_BUT_TEST_ENV)

		}
	}

}
func (uc *UserClientsService) readMessage(clientKey string, rawMsg []byte) error {
	uc.LoggerService.Log(clientKey, ">> ", string(rawMsg))
	msg, unmarshalErr := UnmarshalToMessage(rawMsg)
	if unmarshalErr != nil {
		return fmt.Errorf("error unmarshalling message: %s", unmarshalErr)
	}
	if authErr := uc.AuthService.ValidateAuthInMessage(msg); authErr != nil {
		return fmt.Errorf("error validating auth in message: %s", authErr)
	}
	uc.MsgService.HandleMessage(msg)
	uc.BroadcastMessage(msg)
}

func (uc *UserClientsService) BroadcastMessage(message *Message) {
	msgCopy := *message
	msgCopy.PrivateKey = ""
	subbedClientKeys := uc.SubService.GetClientKeysSubbedToTopic(msgCopy.Topic)
	for _, clientKey := range subbedClientKeys.Flatten() {
		client, err := uc.GetClient(clientKey)
		if err != nil {
			uc.LoggerService.LogRed(ENV_SERVER, fmt.Sprintf("error getting client from key: %s", err), ALL_BUT_TEST_ENV)
			continue
		}
		client.InChannel() <- &msgCopy
	}
}

func (uc *UserClientsService) DirectMessage(message *Message,
	clientKey
string) error{
if message.Topic != "directMessage" && message.Topic != ""{
return fmt.Errorf("direct messages expected to not have a topic, given %s", message.Topic)
}
msgCopy := *message
msgCopy.Topic = "directMessage"
client, err := uc.GetClient(clientKey)
if err != nil{
return err
}
client.InChannel() <- &msgCopy
return nil
}

func (uc *UserClientsService) CleanupClient(userClient *UserClient) {
	_ = uc.AuthService.RemoveClient(userClient.PublicKey())
}
