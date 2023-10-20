package server

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserClientsManager", func() {
	var ucm *UserClientsManager
	var clientKey string
	BeforeEach(func() {
		ucm = NewUserClientsManager()
		clientKey = "some-client-key"
	})
	Describe("::AddClient", func() {
		It("adds the client key to the set", func() {
			err := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.clientKeys.Size()).To(Equal(1))
		})
		It("adds the channel keyed by the client", func() {
			err := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(err).ToNot(HaveOccurred())
			_, ok := ucm.channelByClientKey[clientKey]
			Expect(ok).To(BeTrue())
		})
		When("the client already exists", func() {
			It("returns an error", func() {
				err := ucm.AddClient(clientKey, make(chan *Prompt))
				Expect(err).ToNot(HaveOccurred())
				err2 := ucm.AddClient(clientKey, make(chan *Prompt))
				Expect(err2).To(HaveOccurred())
				Expect(err2).To(Equal(fmt.Errorf("client %s already exists", clientKey)))
			})
		})
	})
	Describe("::SubscribeClientTo", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = MESSAGE_TOPIC_AUTH
			addClientErr := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(addClientErr).ToNot(HaveOccurred())
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
			Expect(ucm.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(0))
		})
		It("adds the topic to the set of subscribed topics keyed by the client key", func() {
			err := ucm.SubscribeClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
			Expect(ucm.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
		})
		When("the client is already subscribed", func() {
			BeforeEach(func() {
				err := ucm.SubscribeClientTo(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
				Expect(ucm.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
			})
			It("returns an error", func() {
				err := ucm.SubscribeClientTo(clientKey, topic)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client %s already subscribed to topic %d", clientKey, topic)))
			})
		})
	})
	Describe("::UnsubClientFrom", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = MESSAGE_TOPIC_AUTH
			addClientErr := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(addClientErr).ToNot(HaveOccurred())
			err := ucm.SubscribeClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
			Expect(ucm.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
		})
		It("removes the topic from the clients subscribed topic", func() {
			err := ucm.UnsubClientFrom(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
			Expect(ucm.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(0))
		})
		When("the client is not subscribed to the topic", func() {
			BeforeEach(func() {
				err := ucm.UnsubClientFrom(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
				Expect(ucm.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(0))
			})
			It("returns an error", func() {
				err := ucm.UnsubClientFrom(clientKey, topic)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client %s is not subscribed to %d", clientKey, topic)))
			})
		})
	})
	Describe("::UnsubClientFromAll", func() {
		var topicA MessageTopic
		var topicB MessageTopic
		BeforeEach(func() {
			topicA = MESSAGE_TOPIC_AUTH
			topicB = MESSAGE_TOPIC_MATCHMAKING
			addClientErr := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(addClientErr).ToNot(HaveOccurred())
			err := ucm.SubscribeClientTo(clientKey, topicA)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
			Expect(ucm.GetClientKeysSubscribedToTopic(topicA).Size()).To(Equal(1))
			err = ucm.SubscribeClientTo(clientKey, topicB)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(2))
			Expect(ucm.GetClientKeysSubscribedToTopic(topicB).Size()).To(Equal(1))
		})
		It("removes both topics from the client's subscriptions", func() {
			ucm.UnsubClientFromAll(clientKey)
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
			Expect(ucm.GetClientKeysSubscribedToTopic(topicA).Size()).To(Equal(0))
			Expect(ucm.GetClientKeysSubscribedToTopic(topicB).Size()).To(Equal(0))
		})
	})
	Describe("::GetClientKeysSubscribedToTopic", func() {
		When("the topic has never been subscribed to", func() {
			BeforeEach(func() {
				ucm = NewUserClientsManager()
			})
			It("initializes an empty set and returns it", func() {
				subbedClientKeys := ucm.GetClientKeysSubscribedToTopic(MESSAGE_TOPIC_AUTH)
				Expect(*subbedClientKeys).To(BeAssignableToTypeOf(set.Set[string]{}))
				Expect(subbedClientKeys.Size()).To(Equal(0))
			})
		})
	})
	Describe("::GetSubscribedTopics", func() {
		When("the client never subscribed to any topics", func() {
			BeforeEach(func() {
				ucm = NewUserClientsManager()
				err := ucm.AddClient(clientKey, make(chan *Prompt))
				Expect(err).ToNot(HaveOccurred())
			})
			It("initializes an empty set and returns it", func() {
				subbedTopics := ucm.GetSubscribedTopics(clientKey)
				Expect(*subbedTopics).To(BeAssignableToTypeOf(set.Set[MessageTopic]{}))
				Expect(subbedTopics.Size()).To(Equal(0))
			})
		})
	})
	Describe("::RemoveClient", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = MESSAGE_TOPIC_AUTH
			addClientErr := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(addClientErr).ToNot(HaveOccurred())
			subClientErr := ucm.SubscribeClientTo(clientKey, topic)
			Expect(subClientErr).ToNot(HaveOccurred())
			Expect(ucm.clientKeys.Has(clientKey)).To(BeTrue())
			Expect(ucm.channelByClientKey[clientKey]).ToNot(BeNil())
			Expect(ucm.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
			Expect(ucm.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
		})
		It("removes the client key from the set", func() {
			err := ucm.RemoveClient(clientKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.clientKeys.Size()).To(Equal(0))
		})
		It("removes the channel keyed by the client", func() {
			err := ucm.RemoveClient(clientKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.channelByClientKey[clientKey]).To(BeNil())
		})
		It("unsubs client from all topics", func() {
			err := ucm.RemoveClient(clientKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.subscriberClientKeysByTopic[topic].Size()).To(Equal(0))
			Expect(ucm.subscribedTopicsByClientKey[clientKey].Size()).To(Equal(0))
		})
		When("the player was never added", func() {
			BeforeEach(func() {
				ucm = NewUserClientsManager()
			})
			It("returns an error", func() {
				err := ucm.RemoveClient(clientKey)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client %s isn't in the clientKey set", clientKey)))
			})
		})
	})
})
