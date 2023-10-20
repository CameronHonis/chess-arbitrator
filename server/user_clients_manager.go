package server

import (
	"fmt"
	. "github.com/CameronHonis/chess-arbitrator/set"
)

type UserClientsManager struct {
	clientKeys                  *Set[string]
	channelByClientKey          map[string]chan *Prompt
	subscriberClientKeysByTopic []*Set[string]
	subscribedTopicsByClientKey map[string]*Set[MessageTopic]
}

func (ucm *UserClientsManager) AddClient(clientKey string, ch chan *Prompt) error {
	if ucm.clientKeys.Has(clientKey) {
		return fmt.Errorf("client %s already exists", clientKey)
	}
	ucm.clientKeys.Add(clientKey)
	ucm.channelByClientKey[clientKey] = ch
	return nil
}

func (ucm *UserClientsManager) RemoveClient(clientKey string) error {
	if !ucm.clientKeys.Has(clientKey) {
		return fmt.Errorf("client %s isn't in the clientKey set", clientKey)
	}
	ucm.clientKeys.Remove(clientKey)
	delete(ucm.channelByClientKey, clientKey)
	ucm.UnsubClientFromAll(clientKey)
	return nil
}

func (ucm *UserClientsManager) SubscribeClientTo(clientKey string, topic MessageTopic) {
	subbedTopics := ucm.GetSubscribedTopics(clientKey)
	subbedTopics.Add(topic)
}

func (ucm *UserClientsManager) UnsubClientFrom(clientKey string, topic MessageTopic) {
	subbedTopics := ucm.GetSubscribedTopics(clientKey)
	subbedTopics.Remove(topic)
}

func (ucm *UserClientsManager) UnsubClientFromAll(clientKey string) {
	for _, topic := range ucm.subscribedTopicsByClientKey[clientKey].Flatten() {
		ucm.subscriberClientKeysByTopic[topic].Remove(clientKey)
	}
	ucm.subscribedTopicsByClientKey[clientKey] = EmptySet[MessageTopic]()
}

func (ucm *UserClientsManager) GetSubscribedTopics(clientKey string) *Set[MessageTopic] {
	subbedTopics, ok := ucm.subscribedTopicsByClientKey[clientKey]
	if !ok {
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

func (ucm *UserClientsManager) GetChannelByClientKey(clientKey string) chan *Prompt {
	return ucm.channelByClientKey[clientKey]
}

func NewUserClientsManager() *UserClientsManager {
	return &UserClientsManager{
		clientKeys:                  EmptySet[string](),
		channelByClientKey:          make(map[string]chan *Prompt),
		subscriberClientKeysByTopic: make([]*Set[string], 50),
		subscribedTopicsByClientKey: make(map[string]*Set[MessageTopic]),
	}
}
