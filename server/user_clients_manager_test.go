package server

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserClientsManager", func() {
	var client *UserClient
	var clientKey string
	BeforeEach(func() {
		userClientsManager = nil
		NewUserClientsManager()
		client = NewUserClient(nil, nil, func(client *UserClient) {})
		client.publicKey = "some-public-key"
		clientKey = client.publicKey
	})
	Describe("::AddClient", func() {
		When("the client hasn't been added", func() {
			It("adds the client to the map", func() {
				err := userClientsManager.AddClient(client)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(userClientsManager.clientByPublicKey)).To(Equal(1))
				client, ok := userClientsManager.clientByPublicKey[clientKey]
				Expect(ok).To(BeTrue())
				Expect(*client).To(BeAssignableToTypeOf(UserClient{}))
				Expect(client.publicKey).To(Equal(clientKey))
				Expect(client.inChannel).ToNot(BeNil())
				Expect(client.outChannel).ToNot(BeNil())
			})
		})
		When("the client already exists", func() {
			BeforeEach(func() {
				err := userClientsManager.AddClient(client)
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns an error", func() {
				err2 := userClientsManager.AddClient(client)
				Expect(err2).To(HaveOccurred())
				Expect(err2).To(Equal(fmt.Errorf("client with key %s already exists", clientKey)))
			})
		})
	})
	Describe("::SubscribeClientTo", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = MESSAGE_TOPIC_AUTH
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(0))
		})
		It("adds the topic to the set of subscribed topics keyed by the client key", func() {
			err := userClientsManager.SubscribeClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
		})
		When("the client is already subscribed", func() {
			BeforeEach(func() {
				err := userClientsManager.SubscribeClientTo(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
				Expect(userClientsManager.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
			})
			It("returns an error", func() {
				err := userClientsManager.SubscribeClientTo(clientKey, topic)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client %s already subscribed to topic %d", clientKey, topic)))
			})
		})
	})
	Describe("::UnsubClientFrom", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = MESSAGE_TOPIC_AUTH
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			err := userClientsManager.SubscribeClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
		})
		It("removes the topic from the clients subscribed topic", func() {
			err := userClientsManager.UnsubClientFrom(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(0))
		})
		When("the client is not subscribed to the topic", func() {
			BeforeEach(func() {
				err := userClientsManager.UnsubClientFrom(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
				Expect(userClientsManager.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(0))
			})
			It("returns an error", func() {
				err := userClientsManager.UnsubClientFrom(clientKey, topic)
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
			topicB = MESSAGE_TOPIC_INIT_BOT_MATCH
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			err := userClientsManager.SubscribeClientTo(clientKey, topicA)
			Expect(err).ToNot(HaveOccurred())
			Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topicA).Size()).To(Equal(1))
			err = userClientsManager.SubscribeClientTo(clientKey, topicB)
			Expect(err).ToNot(HaveOccurred())
			Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(2))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topicB).Size()).To(Equal(1))
		})
		It("removes both topics from the client's subscriptions", func() {
			userClientsManager.UnsubClientFromAll(clientKey)
			Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(0))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topicA).Size()).To(Equal(0))
			Expect(userClientsManager.GetClientKeysSubscribedToTopic(topicB).Size()).To(Equal(0))
		})
	})
	Describe("::GetClientKeysSubscribedToTopic", func() {
		When("the topic has never been subscribed to", func() {
			It("initializes an empty set and returns it", func() {
				subbedClientKeys := userClientsManager.GetClientKeysSubscribedToTopic(MESSAGE_TOPIC_AUTH)
				Expect(*subbedClientKeys).To(BeAssignableToTypeOf(set.Set[string]{}))
				Expect(subbedClientKeys.Size()).To(Equal(0))
			})
		})
	})
	Describe("::GetSubscribedTopics", func() {
		When("the client never subscribed to any topics", func() {
			BeforeEach(func() {
				err := userClientsManager.AddClient(client)
				Expect(err).ToNot(HaveOccurred())
			})
			It("initializes an empty set and returns it", func() {
				subbedTopics := userClientsManager.GetSubscribedTopics(clientKey)
				Expect(*subbedTopics).To(BeAssignableToTypeOf(set.Set[MessageTopic]{}))
				Expect(subbedTopics.Size()).To(Equal(0))
			})
		})
	})
	Describe("::RemoveClient", func() {
		var topic MessageTopic
		When("a player was added", func() {
			BeforeEach(func() {
				topic = MESSAGE_TOPIC_AUTH
				addClientErr := userClientsManager.AddClient(client)
				Expect(addClientErr).ToNot(HaveOccurred())
				_, ok := userClientsManager.clientByPublicKey[clientKey]
				Expect(ok).To(BeTrue())

				subClientErr := userClientsManager.SubscribeClientTo(clientKey, topic)
				Expect(subClientErr).ToNot(HaveOccurred())
				Expect(userClientsManager.GetSubscribedTopics(clientKey).Size()).To(Equal(1))
				Expect(userClientsManager.GetClientKeysSubscribedToTopic(topic).Size()).To(Equal(1))
			})
			It("removes the client from the client map", func() {
				err := userClientsManager.RemoveClient(client)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(userClientsManager.clientByPublicKey)).To(Equal(0))
			})
			It("unsubs client from all topics", func() {
				err := userClientsManager.RemoveClient(client)
				Expect(err).ToNot(HaveOccurred())
				Expect(userClientsManager.subscriberClientKeysByTopic[topic].Size()).To(Equal(0))
				Expect(userClientsManager.subscribedTopicsByClientKey[clientKey].Size()).To(Equal(0))
			})
		})
		When("the player was never added", func() {
			It("returns an error", func() {
				err := userClientsManager.RemoveClient(client)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client with key %s isn't an established client", clientKey)))
			})
		})
	})
	Describe("::GetAllOutChannels", func() {
		var otherClient *UserClient
		BeforeEach(func() {
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			otherClient = NewUserClient(nil, nil, func(client *UserClient) {})
			otherClient.publicKey = "other-public-key"
			addClientErr = userClientsManager.AddClient(otherClient)
			Expect(addClientErr).ToNot(HaveOccurred())
		})
		It("returns a slice of all client channels", func() {
			channels := userClientsManager.GetAllOutChannels()
			Expect(channels).To(HaveLen(2))
			Expect(&channels[0]).To(Equal(&client.outChannel))
			Expect(&channels[1]).To(Equal(&otherClient.outChannel))
		})
	})
})
