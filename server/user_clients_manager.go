package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	"github.com/gorilla/websocket"
	"sync"
	"time"
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
	AddNewClient(conn *websocket.Conn) (*UserClient, error)
	AddClient(client *UserClient) error
	RemoveClient(client *UserClient) error
	GetClientFromKey(clientKey string) (*UserClient, error)
	GetAllOutChannels() map[string]chan *Message
	GetInChannelByClientKey(clientKey string) (<-chan *Message, error)
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

func (ucm *UserClientsService) AddNewClient(conn *websocket.Conn) (*UserClient, error) {
	client := NewUserClient(conn, CleanupClient)

	err := ucm.AddClient(client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (ucm *UserClientsService) AddClient(client *UserClient) error {
	ucm.mu.Lock()
	defer ucm.mu.Unlock()
	if _, ok := ucm.clientByPublicKey[client.PublicKey()]; ok {
		client.Kill()
		return fmt.Errorf("client with key %s already exists", client.PublicKey())
	}
	ucm.clientByPublicKey[client.PublicKey()] = client
	return nil
}

func (ucm *UserClientsService) RemoveClient(client *UserClient) error {
	pubKey := client.PublicKey()

	ucm.mu.Lock()
	if _, ok := ucm.clientByPublicKey[pubKey]; !ok {
		ucm.mu.Unlock()
		return fmt.Errorf("client with key %s isn't an established client", client.PublicKey())
	}
	delete(ucm.clientByPublicKey, pubKey)
	ucm.mu.Unlock()

	client.Kill()
	ucm.SubService.UnsubClientFromAll(pubKey)
	return nil
}

func (ucm *UserClientsService) GetClientFromKey(clientKey string) (*UserClient, error) {
	defer ucm.mu.Unlock()
	ucm.mu.Lock()
	client, ok := ucm.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("no client with public key %s", clientKey)
	}
	return client, nil
}

func (ucm *UserClientsService) GetAllOutChannels() map[string]chan *Message {
	defer ucm.mu.Unlock()
	ucm.mu.Lock()
	if len(ucm.clientByPublicKey) == 0 {
		return make(map[string]chan *Message)
	}
	channels := make(map[string]chan *Message)
	for _, client := range ucm.clientByPublicKey {
		channels[client.PublicKey()] = client.OutChannel()
	}
	return channels
}

func (ucm *UserClientsService) GetInChannelByClientKey(clientKey string) (<-chan *Message, error) {
	defer ucm.mu.Unlock()
	ucm.mu.Lock()
	client, ok := ucm.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("client with key %s does not exist", clientKey)
	}
	return client.OutChannel(), nil
}

func (ucm *UserClientsService) listenOnUserClientChannels() {
	for {
		time.Sleep(time.Millisecond * 1)
		for clientKey, channel := range ucm.GetAllOutChannels() {
			select {
			case message := <-channel:
				message.SenderKey = clientKey
				go ucm.MsgService.HandleMessage(message)
			default:
				continue
			}
		}
	}
}

func (ucm *UserClientsService) BroadcastMessage(message *Message) {
	msgCopy := *message
	msgCopy.PrivateKey = ""
	subbedClientKeys := ucm.SubService.GetClientKeysSubbedToTopic(msgCopy.Topic)
	for _, clientKey := range subbedClientKeys.Flatten() {
		client, err := ucm.GetClientFromKey(clientKey)
		if err != nil {
			ucm.LoggerService.LogRed(ENV_SERVER, fmt.Sprintf("error getting client from key: %s", err), ALL_BUT_TEST_ENV)
			continue
		}
		client.InChannel() <- &msgCopy
	}
}

func (ucm *UserClientsService) DirectMessage(message *Message, clientKey string) error {
	if message.Topic != "directMessage" && message.Topic != "" {
		return fmt.Errorf("direct messages expected to not have a topic, given %s", message.Topic)
	}
	msgCopy := *message
	msgCopy.Topic = "directMessage"
	client, err := ucm.GetClientFromKey(clientKey)
	if err != nil {
		return err
	}
	client.InChannel() <- &msgCopy
	return nil
}

func (ucm *UserClientsService) CleanupClient(userClient *UserClient) {
	if ucm.AuthService.GetBotKey() == userClient.PublicKey() {
		am.chessBotKey = ""
	}
}
