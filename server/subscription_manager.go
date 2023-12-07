package server

import (
	"fmt"
	. "github.com/CameronHonis/set"
	"sync"
)

var subscriptionManager *SubscriptionManager

type MessageTopic string
type SubscriptionManagerI interface {
	SubClientTo(clientKey string, topic MessageTopic) error
	UnsubClientFrom(clientKey string, topic MessageTopic) error
	UnsubClientFromAll(clientKey string)
	GetSubbedTopics(clientKey string) *Set[MessageTopic]
	GetClientKeysSubbedToTopic(topic MessageTopic) *Set[string]
}

type SubscriptionManager struct {
	userClientsManager UserClientsManagerI

	subscriberClientKeysByTopic map[MessageTopic]*Set[string]
	subbedTopicsByClientKey     map[string]*Set[MessageTopic]
	mu                          sync.Mutex
}

func GetSubscriptionManager() *SubscriptionManager {
	if subscriptionManager == nil {
		subscriptionManager = &SubscriptionManager{
			userClientsManager:          GetUserClientsManager(),
			subscriberClientKeysByTopic: make(map[MessageTopic]*Set[string]),
			subbedTopicsByClientKey:     make(map[string]*Set[MessageTopic]),
		}
	}
	return subscriptionManager
}

func (sm *SubscriptionManager) SubClientTo(clientKey string, topic MessageTopic) error {
	subbedTopics := sm.GetSubbedTopics(clientKey)
	sm.mu.Lock()
	if subbedTopics.Has(topic) {
		return fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Add(topic)
	sm.mu.Unlock()

	subbedClientKeys := sm.GetClientKeysSubbedToTopic(topic)

	sm.mu.Lock()
	subbedClientKeys.Add(clientKey)
	sm.mu.Unlock()
	return nil
}

func (sm *SubscriptionManager) UnsubClientFrom(clientKey string, topic MessageTopic) error {
	subbedTopics := sm.GetSubbedTopics(clientKey)
	sm.mu.Lock()
	if !subbedTopics.Has(topic) {
		return fmt.Errorf("client %s not subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Remove(topic)
	sm.mu.Unlock()

	subbedClientKeys := sm.GetClientKeysSubbedToTopic(topic)

	sm.mu.Lock()
	subbedClientKeys.Remove(clientKey)
	sm.mu.Unlock()
	return nil
}

func (sm *SubscriptionManager) UnsubClientFromAll(clientKey string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	subbedTopics := sm.subbedTopicsByClientKey[clientKey]
	for _, topic := range subbedTopics.Flatten() {
		sm.subscriberClientKeysByTopic[topic].Remove(clientKey)
	}
	delete(sm.subbedTopicsByClientKey, clientKey)
}

func (sm *SubscriptionManager) GetSubbedTopics(clientKey string) *Set[MessageTopic] {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.subbedTopicsByClientKey[clientKey]
}

func (sm *SubscriptionManager) GetClientKeysSubbedToTopic(topic MessageTopic) *Set[string] {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.subscriberClientKeysByTopic[topic]
}
