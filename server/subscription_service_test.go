package server

import (
	"fmt"
	"github.com/CameronHonis/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionService", func() {
	var client *Client
	var clientKey string
	BeforeEach(func() {
		userClientsManager = nil
		GetUserClientsManager()
		subscriptionManager = nil
		GetSubscriptionManager()
		client = NewClient(nil, func(client *Client) {})
		client.publicKey = "some-public-key"
		clientKey = client.publicKey
	})
	Describe("::SubClientTo", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = "auth"
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			subbedTopics := subscriptionManager.GetSubbedTopics(clientKey)
			Expect(subbedTopics.Size()).To(Equal(0))
			subbedClientKeys := subscriptionManager.GetClientKeysSubbedToTopic(topic)
			Expect(subbedClientKeys.Size()).To(Equal(0))
		})
		It("adds the topic to the set of subscribed topics keyed by the client key", func() {
			err := subscriptionManager.SubClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(1))
			Expect(subscriptionManager.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
		})
		When("the client is already subscribed", func() {
			BeforeEach(func() {
				err := subscriptionManager.SubClientTo(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(1))
				Expect(subscriptionManager.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
			})
			It("returns an error", func() {
				err := subscriptionManager.SubClientTo(clientKey, topic)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)))
			})
		})
	})
	Describe("::UnsubClientFrom", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = "auth"
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			err := subscriptionManager.SubClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(1))
			Expect(subscriptionManager.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
		})
		It("removes the topic from the clients subscribed topic", func() {
			err := subscriptionManager.UnsubClientFrom(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(0))
			Expect(subscriptionManager.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
		})
		When("the client is not subscribed to the topic", func() {
			BeforeEach(func() {
				err := subscriptionManager.UnsubClientFrom(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(0))
				Expect(subscriptionManager.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
			})
			It("returns an error", func() {
				err := subscriptionManager.UnsubClientFrom(clientKey, topic)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client %s not subscribed to topic %s", clientKey, topic)))
			})
		})
	})
	Describe("::UnsubClientFromAll", func() {
		var topicA MessageTopic
		var topicB MessageTopic
		BeforeEach(func() {
			topicA = "auth"
			topicB = "findMatch"
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			err := subscriptionManager.SubClientTo(clientKey, topicA)
			Expect(err).ToNot(HaveOccurred())
			Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(1))
			Expect(subscriptionManager.GetClientKeysSubbedToTopic(topicA).Size()).To(Equal(1))
			err = subscriptionManager.SubClientTo(clientKey, topicB)
			Expect(err).ToNot(HaveOccurred())
			Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(2))
			Expect(subscriptionManager.GetClientKeysSubbedToTopic(topicB).Size()).To(Equal(1))
		})
		It("removes both topics from the client's subscriptions", func() {
			subscriptionManager.UnsubClientFromAll(clientKey)
			Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(0))
			Expect(subscriptionManager.GetClientKeysSubbedToTopic(topicA).Size()).To(Equal(0))
			Expect(subscriptionManager.GetClientKeysSubbedToTopic(topicB).Size()).To(Equal(0))
		})
		When("the client has never subscribed to any topics", func() {
			It("does not remove subscribed topics from other clients", func() {
				subscriptionManager.UnsubClientFromAll("some-random-client-key")
				Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(2))
			})
		})
	})
	Describe("::GetClientKeysSubbedToTopic", func() {
		When("the topic has never been subscribed to", func() {
			It("initializes an empty set and returns it", func() {
				subbedClientKeys := subscriptionManager.GetClientKeysSubbedToTopic("auth")
				Expect(*subbedClientKeys).To(BeAssignableToTypeOf(set.Set[string]{}))
				Expect(subbedClientKeys.Size()).To(Equal(0))
			})
		})
		When("the topic has been subscribed to", func() {
			BeforeEach(func() {
				subErr := subscriptionManager.SubClientTo(clientKey, "topicA")
				Expect(subErr).ToNot(HaveOccurred())
				subErr = subscriptionManager.SubClientTo(clientKey, "topicB")
				Expect(subErr).ToNot(HaveOccurred())
			})
			It("returns the set of client keys subscribed to the topic", func() {
				subbedClientKeys := subscriptionManager.GetClientKeysSubbedToTopic("topicA")
				Expect(subbedClientKeys.Size()).To(Equal(1))
				Expect(subbedClientKeys.Has(clientKey)).To(BeTrue())
			})
		})
	})
	Describe("::GetSubbedTopics", func() {
		BeforeEach(func() {
			userClientsManager = nil
			GetUserClientsManager()
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
		})
		When("the client never subscribed to any topics", func() {
			It("initializes an empty set and returns it", func() {
				subbedTopics := subscriptionManager.GetSubbedTopics(clientKey)
				Expect(*subbedTopics).To(BeAssignableToTypeOf(set.Set[MessageTopic]{}))
				Expect(subbedTopics.Size()).To(Equal(0))
			})
		})
		When("the client has subscribed to topics", func() {
			BeforeEach(func() {
				subErr := subscriptionManager.SubClientTo(clientKey, "topicA")
				Expect(subErr).ToNot(HaveOccurred())
				subErr = subscriptionManager.SubClientTo(clientKey, "topicB")
				Expect(subErr).ToNot(HaveOccurred())
			})
			It("returns the set of topics the client is subscribed to", func() {
				subbedTopics := subscriptionManager.GetSubbedTopics(clientKey)
				Expect(subbedTopics.Size()).To(Equal(2))
				Expect(subbedTopics.Has("topicA")).To(BeTrue())
				Expect(subbedTopics.Has("topicB")).To(BeTrue())
			})
		})
	})
})
