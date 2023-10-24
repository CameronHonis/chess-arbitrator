package server

import (
	"fmt"
	. "github.com/CameronHonis/chess-arbitrator/set"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var userClientsManager *UserClientsManager

type UserClientsManager struct {
	interactMutex               *sync.Mutex
	clientByPublicKey           map[string]*UserClient
	subscriberClientKeysByTopic []*Set[string]
	subscribedTopicsByClientKey map[string]*Set[MessageTopic]
}

func (ucm *UserClientsManager) AddNewClient(conn *websocket.Conn) (*UserClient, error) {
	client := NewUserClient(conn, func(client *UserClient) {})

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
	ucm.UnsubClientFromAll(client.PublicKey())
	return nil
}

func (ucm *UserClientsManager) SubscribeClientTo(clientKey string, topic MessageTopic) error {
	subbedTopics := ucm.GetSubscribedTopics(clientKey)
	ucm.interactMutex.Lock()
	if subbedTopics.Has(topic) {
		return fmt.Errorf("client %s already subscribed to topic %d", clientKey, topic)
	}
	subbedTopics.Add(topic)
	ucm.interactMutex.Unlock()

	subbedClientKeys := ucm.GetClientKeysSubscribedToTopic(topic)

	ucm.interactMutex.Lock()
	subbedClientKeys.Add(clientKey)
	ucm.interactMutex.Unlock()
	return nil
}

func (ucm *UserClientsManager) UnsubClientFrom(clientKey string, topic MessageTopic) error {
	defer ucm.interactMutex.Unlock()
	subbedTopics := ucm.GetSubscribedTopics(clientKey)

	ucm.interactMutex.Lock()
	if !subbedTopics.Has(topic) {
		return fmt.Errorf("client %s is not subscribed to %d", clientKey, topic)
	}
	subbedTopics.Remove(topic)
	ucm.interactMutex.Unlock()

	subbedClientKeys := ucm.GetClientKeysSubscribedToTopic(topic)

	ucm.interactMutex.Lock()
	subbedClientKeys.Remove(clientKey)
	return nil
}

func (ucm *UserClientsManager) UnsubClientFromAll(clientKey string) {
	defer ucm.interactMutex.Unlock()
	subbedTopics := ucm.GetSubscribedTopics(clientKey)
	ucm.interactMutex.Lock()
	for _, topic := range subbedTopics.Flatten() {
		ucm.subscriberClientKeysByTopic[topic].Remove(clientKey)
	}
	ucm.subscribedTopicsByClientKey[clientKey] = EmptySet[MessageTopic]()
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

func (ucm *UserClientsManager) GetSubscribedTopics(clientKey string) *Set[MessageTopic] {
	defer ucm.interactMutex.Unlock()
	ucm.interactMutex.Lock()
	subbedTopics, ok := ucm.subscribedTopicsByClientKey[clientKey]
	if !ok {
		subbedTopics = EmptySet[MessageTopic]()
		ucm.subscribedTopicsByClientKey[clientKey] = subbedTopics
	}
	return subbedTopics
}

func (ucm *UserClientsManager) GetClientKeysSubscribedToTopic(topic MessageTopic) *Set[string] {
	defer ucm.interactMutex.Unlock()
	ucm.interactMutex.Lock()
	subbedClients := ucm.subscriberClientKeysByTopic[topic]
	if subbedClients == nil {
		subbedClients = EmptySet[string]()
		ucm.subscriberClientKeysByTopic[topic] = subbedClients
	}
	return subbedClients
}

func (ucm *UserClientsManager) GetAllOutChannels() []chan *Prompt {
	defer ucm.interactMutex.Unlock()
	ucm.interactMutex.Lock()
	if len(ucm.clientByPublicKey) == 0 {
		return []chan *Prompt{}
	}
	channels := make([]chan *Prompt, len(ucm.clientByPublicKey))
	channelIdx := 0
	for _, client := range ucm.clientByPublicKey {
		channels[channelIdx] = client.OutChannel()
		channelIdx++
	}
	return channels
}

func (ucm *UserClientsManager) GetInChannelByClientKey(clientKey string) (<-chan *Prompt, error) {
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
		for _, channel := range ucm.GetAllOutChannels() {
			select {
			case prompt := <-channel:
				go handlePrompt(prompt)
			default:
				continue
			}
		}
	}
}

func GetUserClientsManager() *UserClientsManager {
	if userClientsManager != nil {
		return userClientsManager
	}
	ucm := UserClientsManager{
		interactMutex:               &sync.Mutex{},
		clientByPublicKey:           make(map[string]*UserClient),
		subscriberClientKeysByTopic: make([]*Set[string], 50),
		subscribedTopicsByClientKey: make(map[string]*Set[MessageTopic]),
	}
	go ucm.listenOnUserClientChannels()
	userClientsManager = &ucm
	return &ucm
}
