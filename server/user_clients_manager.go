package server

import (
	"fmt"
	. "github.com/CameronHonis/chess-arbitrator/set"
	"github.com/gorilla/websocket"
	"sync"
)

var userClientsManager *UserClientsManager

type UserClientsManager struct {
	stdoutMutex                 *sync.Mutex
	clientByPublicKey           map[string]*UserClient
	subscriberClientKeysByTopic []*Set[string]
	subscribedTopicsByClientKey map[string]*Set[MessageTopic]
}

func (ucm *UserClientsManager) AddNewClient(conn *websocket.Conn) (*UserClient, error) {
	clientChannel := make(chan *Prompt)
	client := NewUserClient(ucm.stdoutMutex, clientChannel, conn, func(client *UserClient) {})

	err := ucm.AddClient(client)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (ucm *UserClientsManager) AddClient(client *UserClient) error {
	if _, ok := ucm.clientByPublicKey[client.PublicKey()]; ok {
		return fmt.Errorf("client with key %s already exists", client.PublicKey())
	}
	ucm.clientByPublicKey[client.PublicKey()] = client
	return nil
}

func (ucm *UserClientsManager) RemoveClient(client *UserClient) error {
	if _, ok := ucm.clientByPublicKey[client.PublicKey()]; !ok {
		return fmt.Errorf("client with key %s isn't an established client", client.PublicKey())
	}
	delete(ucm.clientByPublicKey, client.PublicKey())
	ucm.UnsubClientFromAll(client.PublicKey())
	return nil
}

func (ucm *UserClientsManager) SubscribeClientTo(clientKey string, topic MessageTopic) error {
	subbedTopics := ucm.GetSubscribedTopics(clientKey)
	if subbedTopics.Has(topic) {
		return fmt.Errorf("client %s already subscribed to topic %d", clientKey, topic)
	}
	subbedTopics.Add(topic)
	subbedClientKeys := ucm.GetClientKeysSubscribedToTopic(topic)
	subbedClientKeys.Add(clientKey)
	return nil
}

func (ucm *UserClientsManager) UnsubClientFrom(clientKey string, topic MessageTopic) error {
	subbedTopics := ucm.GetSubscribedTopics(clientKey)
	if !subbedTopics.Has(topic) {
		return fmt.Errorf("client %s is not subscribed to %d", clientKey, topic)
	}
	subbedTopics.Remove(topic)
	subbedClientKeys := ucm.GetClientKeysSubscribedToTopic(topic)
	subbedClientKeys.Remove(clientKey)
	return nil
}

func (ucm *UserClientsManager) UnsubClientFromAll(clientKey string) {
	subbedTopics := ucm.GetSubscribedTopics(clientKey)
	for _, topic := range subbedTopics.Flatten() {
		ucm.subscriberClientKeysByTopic[topic].Remove(clientKey)
	}
	ucm.subscribedTopicsByClientKey[clientKey] = EmptySet[MessageTopic]()
}

func (ucm *UserClientsManager) GetClientFromKey(clientKey string) (*UserClient, error) {
	client, ok := ucm.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("no client with public key %s", clientKey)
	}
	return client, nil
}

func (ucm *UserClientsManager) GetSubscribedTopics(clientKey string) *Set[MessageTopic] {
	subbedTopics, ok := ucm.subscribedTopicsByClientKey[clientKey]
	if !ok {
		subbedTopics = EmptySet[MessageTopic]()
		ucm.subscribedTopicsByClientKey[clientKey] = subbedTopics
	}
	return subbedTopics
}

func (ucm *UserClientsManager) GetClientKeysSubscribedToTopic(topic MessageTopic) *Set[string] {
	subbedClients := ucm.subscriberClientKeysByTopic[topic]
	if subbedClients == nil {
		subbedClients = EmptySet[string]()
		ucm.subscriberClientKeysByTopic[topic] = subbedClients
	}
	return subbedClients
}

func (ucm *UserClientsManager) GetAllChannels() []chan *Prompt {
	channels := make([]chan *Prompt, len(ucm.clientByPublicKey))
	for _, client := range ucm.clientByPublicKey {
		channels = append(channels, client.ServerChannel())
	}
	return channels
}

func (ucm *UserClientsManager) GetChannelByClientKey(clientKey string) (chan *Prompt, error) {
	client, ok := ucm.clientByPublicKey[clientKey]
	if !ok {
		return nil, fmt.Errorf("client with key %s does not exist", clientKey)
	}
	return client.ServerChannel(), nil
}
func (ucm *UserClientsManager) listenOnUserClientChannels() {
	for {
		for _, channel := range ucm.GetAllChannels() {
			select {
			case prompt := <-channel:
				handlePrompt(prompt)
			default:
				continue
			}
		}
	}
}

func NewUserClientsManager() (*UserClientsManager, error) {
	if userClientsManager != nil {
		return nil, fmt.Errorf("singleton UserClientsManager already instantiated")
	}
	ucm := UserClientsManager{
		stdoutMutex:                 &sync.Mutex{},
		clientByPublicKey:           make(map[string]*UserClient),
		subscriberClientKeysByTopic: make([]*Set[string], 50),
		subscribedTopicsByClientKey: make(map[string]*Set[MessageTopic]),
	}
	go ucm.listenOnUserClientChannels()
	userClientsManager = &ucm
	return &ucm, nil
}
