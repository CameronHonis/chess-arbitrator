package server

import (
	"fmt"
	. "github.com/CameronHonis/log"
	. "github.com/CameronHonis/set"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var userClientsManager *UserClientsManager

type UserClientsManagerI interface {
	AddNewClient(conn *websocket.Conn) (*UserClient, error)
	AddClient(client *UserClient) error
	RemoveClient(client *UserClient) error
	GetClientFromKey(clientKey string) (*UserClient, error)
	GetAllOutChannels() map[string]chan *Message
	GetInChannelByClientKey(clientKey string) (<-chan *Message, error)
	BroadcastMessage(message *Message)
	DirectMessage(message *Message, clientKey string) error
}

type UserClientsManager struct {
	messageHandler      MessageHandlerI
	logManager          LogManagerI
	subscriptionManager SubscriptionManagerI

	interactMutex               sync.Mutex
	clientByPublicKey           map[string]*UserClient
	subscriberClientKeysByTopic map[MessageTopic]*Set[string]
	subscribedTopicsByClientKey map[string]*Set[MessageTopic]
}

func GetUserClientsManager() *UserClientsManager {
	if userClientsManager != nil {
		return userClientsManager
	}
	ucm := UserClientsManager{
		messageHandler:              GetMessageHandler(),
		logManager:                  GetLogManager(),
		interactMutex:               sync.Mutex{},
		clientByPublicKey:           make(map[string]*UserClient),
		subscriberClientKeysByTopic: make(map[MessageTopic]*Set[string], 50),
		subscribedTopicsByClientKey: make(map[string]*Set[MessageTopic]),
	}
	go ucm.listenOnUserClientChannels()
	userClientsManager = &ucm
	return &ucm
}

func (ucm *UserClientsManager) AddNewClient(conn *websocket.Conn) (*UserClient, error) {
	client := NewUserClient(conn, CleanupClient)

	err := ucm.AddClient(client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (ucm *UserClientsManager) AddClient(client *UserClient) error {
	ucm.interactMutex.Lock()
	defer ucm.interactMutex.Unlock()
	if _, ok := ucm.clientByPublicKey[client.PublicKey()]; ok {
		client.Kill()
		return fmt.Errorf("client with key %s already exists", client.PublicKey())
	}
	ucm.clientByPublicKey[client.PublicKey()] = client
	return nil
}

func (ucm *UserClientsManager) RemoveClient(client *UserClient) error {
	ucm.interactMutex.Lock()
	if _, ok := ucm.clientByPublicKey[client.PublicKey()]; !ok {
		ucm.interactMutex.Unlock()
		return fmt.Errorf("client with key %s isn't an established client", client.PublicKey())
	}
	delete(ucm.clientByPublicKey, client.PublicKey())
	ucm.interactMutex.Unlock()

	client.Kill()
	ucm.subscriptionManager.UnsubClientFromAll(client.PublicKey())
	return nil
}

func (ucm *UserClientsManager) GetClientFromKey(clientKey string) (*UserClient, error) {
	defer ucm.interactMutex.Unlock()
	ucm.interactMutex.Lock()
	client, ok := ucm.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("no client with public key %s", clientKey)
	}
	return client, nil
}

func (ucm *UserClientsManager) GetAllOutChannels() map[string]chan *Message {
	defer ucm.interactMutex.Unlock()
	ucm.interactMutex.Lock()
	if len(ucm.clientByPublicKey) == 0 {
		return make(map[string]chan *Message)
	}
	channels := make(map[string]chan *Message)
	for _, client := range ucm.clientByPublicKey {
		channels[client.PublicKey()] = client.OutChannel()
	}
	return channels
}

func (ucm *UserClientsManager) GetInChannelByClientKey(clientKey string) (<-chan *Message, error) {
	defer ucm.interactMutex.Unlock()
	ucm.interactMutex.Lock()
	client, ok := ucm.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("client with key %s does not exist", clientKey)
	}
	return client.OutChannel(), nil
}

func (ucm *UserClientsManager) listenOnUserClientChannels() {
	for {
		time.Sleep(time.Millisecond * 1)
		for clientKey, channel := range ucm.GetAllOutChannels() {
			select {
			case message := <-channel:
				message.SenderKey = clientKey
				go ucm.messageHandler.HandleMessage(message)
			default:
				continue
			}
		}
	}
}

func (ucm *UserClientsManager) BroadcastMessage(message *Message) {
	msgCopy := *message
	msgCopy.PrivateKey = ""
	subbedClientKeys := ucm.subscriptionManager.GetClientKeysSubbedToTopic(msgCopy.Topic)
	for _, clientKey := range subbedClientKeys.Flatten() {
		client, err := ucm.GetClientFromKey(clientKey)
		if err != nil {
			ucm.logManager.LogRed(ENV_SERVER, fmt.Sprintf("error getting client from key: %s", err), ALL_BUT_TEST_ENV)
			continue
		}
		client.InChannel() <- &msgCopy
	}
}

func (ucm *UserClientsManager) DirectMessage(message *Message, clientKey string) error {
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

func CleanupClient(userClient *UserClient) {
	am := GetAuthManager()
	if am.chessBotKey == userClient.PublicKey() {
		am.chessBotKey = ""
	}
}
