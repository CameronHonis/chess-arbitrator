package server

import (
	"fmt"
	"github.com/CameronHonis/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionService", func() {
	var client *Client
	var clientKey Key
	var subService *SubscriptionService
	BeforeEach(func() {
		client = NewClient(nil, func(client *Client) {})
		client.publicKey = "some-public-key"
		clientKey = client.PublicKey()
		subService = NewSubscriptionService(NewSubscriptionConfig())
	})
	Describe("::SubClientTo", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = "auth"
			subbedTopics := subService.GetSubbedTopics(clientKey)
			Expect(subbedTopics.Size()).To(Equal(0))
			subbedClientKeys := subService.GetClientKeysSubbedToTopic(topic)
			Expect(subbedClientKeys.Size()).To(Equal(0))
		})
		It("adds the topic to the set of subscribed topics keyed by the client key", func() {
			err := subService.SubClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(1))
			Expect(subService.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
		})
		When("the client is already subscribed", func() {
			BeforeEach(func() {
				err := subService.SubClientTo(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(1))
				Expect(subService.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
			})
			It("returns an error", func() {
				err := subService.SubClientTo(clientKey, topic)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)))
			})
		})
	})
	Describe("::UnsubClientFrom", func() {
		var topic MessageTopic
		BeforeEach(func() {
			topic = "auth"
			err := subService.SubClientTo(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(1))
			Expect(subService.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
		})
		It("removes the topic from the clients subscribed topic", func() {
			err := subService.UnsubClientFrom(clientKey, topic)
			Expect(err).ToNot(HaveOccurred())
			Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(0))
			Expect(subService.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
		})
		When("the client is not subscribed to the topic", func() {
			BeforeEach(func() {
				err := subService.UnsubClientFrom(clientKey, topic)
				Expect(err).ToNot(HaveOccurred())
				Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(0))
				Expect(subService.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
			})
			It("returns an error", func() {
				err := subService.UnsubClientFrom(clientKey, topic)
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
			err := subService.SubClientTo(clientKey, topicA)
			Expect(err).ToNot(HaveOccurred())
			Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(1))
			Expect(subService.GetClientKeysSubbedToTopic(topicA).Size()).To(Equal(1))
			err = subService.SubClientTo(clientKey, topicB)
			Expect(err).ToNot(HaveOccurred())
			Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(2))
			Expect(subService.GetClientKeysSubbedToTopic(topicB).Size()).To(Equal(1))
		})
		It("removes both topics from the client's subscriptions", func() {
			subService.UnsubClientFromAll(clientKey)
			Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(0))
			Expect(subService.GetClientKeysSubbedToTopic(topicA).Size()).To(Equal(0))
			Expect(subService.GetClientKeysSubbedToTopic(topicB).Size()).To(Equal(0))
		})
		When("the client has never subscribed to any topics", func() {
			It("does not remove subscribed topics from other clients", func() {
				subService.UnsubClientFromAll("some-random-client-key")
				Expect(subService.GetSubbedTopics(clientKey).Size()).To(Equal(2))
			})
		})
	})
	Describe("::GetClientKeysSubbedToTopic", func() {
		When("the topic has never been subscribed to", func() {
			It("initializes an empty set and returns it", func() {
				subbedClientKeys := subService.GetClientKeysSubbedToTopic("auth")
				Expect(*subbedClientKeys).To(BeAssignableToTypeOf(set.Set[string]{}))
				Expect(subbedClientKeys.Size()).To(Equal(0))
			})
		})
		When("the topic has been subscribed to", func() {
			BeforeEach(func() {
				subErr := subService.SubClientTo(clientKey, "topicA")
				Expect(subErr).ToNot(HaveOccurred())
				subErr = subService.SubClientTo(clientKey, "topicB")
				Expect(subErr).ToNot(HaveOccurred())
			})
			It("returns the set of client keys subscribed to the topic", func() {
				subbedClientKeys := subService.GetClientKeysSubbedToTopic("topicA")
				Expect(subbedClientKeys.Size()).To(Equal(1))
				Expect(subbedClientKeys.Has(clientKey)).To(BeTrue())
			})
		})
	})
	Describe("::GetSubbedTopics", func() {
		When("the client never subscribed to any topics", func() {
			It("initializes an empty set and returns it", func() {
				subbedTopics := subService.GetSubbedTopics(clientKey)
				Expect(subbedTopics).To(BeAssignableToTypeOf(set.Set[MessageTopic]{}))
				Expect(subbedTopics.Size()).To(Equal(0))
			})
		})
		When("the client has subscribed to topics", func() {
			BeforeEach(func() {
				subErr := subService.SubClientTo(clientKey, "topicA")
				Expect(subErr).ToNot(HaveOccurred())
				subErr = subService.SubClientTo(clientKey, "topicB")
				Expect(subErr).ToNot(HaveOccurred())
			})
			It("returns the set of topics the client is subscribed to", func() {
				subbedTopics := subService.GetSubbedTopics(clientKey)
				Expect(subbedTopics.Size()).To(Equal(2))
				Expect(subbedTopics.Has("topicA")).To(BeTrue())
				Expect(subbedTopics.Has("topicB")).To(BeTrue())
			})
		})
	})
})
