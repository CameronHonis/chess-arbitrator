package sub_service

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"github.com/CameronHonis/set"
	"sync"
)

type SubscriptionServiceI interface {
	service.ServiceI
	SubClient(clientKey models.Key, topic models.MessageTopic) error
	UnsubClient(clientKey models.Key, topic models.MessageTopic) error
	UnsubClientFromAll(clientKey models.Key)
	SubbedTopics(clientKey models.Key) *set.Set[models.MessageTopic]
	ClientKeysSubbedToTopic(topic models.MessageTopic) *set.Set[models.Key]
}

type SubscriptionService struct {
	service.Service
	__dependencies__ marker.Marker
	AuthService      auth.AuthenticationServiceI

	__state__               marker.Marker
	subbedClientKeysByTopic map[models.MessageTopic]*set.Set[models.Key]
	subbedTopicsByClientKey map[models.Key]*set.Set[models.MessageTopic]
	mu                      sync.Mutex
}

func NewSubscriptionService(config *SubscriptionServiceConfig) *SubscriptionService {
	subService := &SubscriptionService{
		subbedClientKeysByTopic: make(map[models.MessageTopic]*set.Set[models.Key]),
		subbedTopicsByClientKey: make(map[models.Key]*set.Set[models.MessageTopic]),
	}
	subService.Service = *service.NewService(subService, config)
	return subService
}
func (s *SubscriptionService) SubClient(clientKey models.Key, topic models.MessageTopic) error {
	authErr := s.AuthService.ValidateClientForTopic(clientKey, topic)
	if authErr != nil {
		go s.Dispatch(NewSubFailedEvent(clientKey, topic, authErr.Error()))
		return authErr
	}
	subbedTopics := s.SubbedTopics(clientKey)
	s.mu.Lock()
	if subbedTopics.Has(topic) {
		go s.Dispatch(NewSubFailedEvent(clientKey, topic, "already subscribed"))
		return fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Add(topic)
	s.mu.Unlock()

	subbedClientKeys := s.ClientKeysSubbedToTopic(topic)

	s.mu.Lock()
	subbedClientKeys.Add(clientKey)
	s.mu.Unlock()

	go s.Dispatch(NewSubSuccessEvent(clientKey, topic))
	return nil
}

func (s *SubscriptionService) UnsubClient(clientKey models.Key, topic models.MessageTopic) error {
	subbedTopics := s.SubbedTopics(clientKey)
	s.mu.Lock()
	if !subbedTopics.Has(topic) {
		return fmt.Errorf("client %s not subscribed to topic %s", clientKey, topic)
	}
	subbedTopics.Remove(topic)
	s.mu.Unlock()

	subbedClientKeys := s.ClientKeysSubbedToTopic(topic)

	s.mu.Lock()
	subbedClientKeys.Remove(clientKey)
	s.mu.Unlock()
	return nil
}

func (s *SubscriptionService) UnsubClientFromAll(clientKey models.Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subbedTopics, ok := s.subbedTopicsByClientKey[clientKey]
	if !ok {
		return
	}
	for _, topic := range subbedTopics.Flatten() {
		s.subbedClientKeysByTopic[topic].Remove(clientKey)
	}
	delete(s.subbedTopicsByClientKey, clientKey)
}

func (s *SubscriptionService) SubbedTopics(clientKey models.Key) *set.Set[models.MessageTopic] {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.subbedTopicsByClientKey[clientKey]
	if !ok {
		s.subbedTopicsByClientKey[clientKey] = set.EmptySet[models.MessageTopic]()
	}
	return s.subbedTopicsByClientKey[clientKey]
}

func (s *SubscriptionService) ClientKeysSubbedToTopic(topic models.MessageTopic) *set.Set[models.Key] {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.subbedClientKeysByTopic[topic]
	if !ok {
		s.subbedClientKeysByTopic[topic] = set.EmptySet[models.Key]()
	}
	return s.subbedClientKeysByTopic[topic]
}
