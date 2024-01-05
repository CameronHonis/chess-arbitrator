package server_test

import (
	"github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/CameronHonis/chess-arbitrator/server/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionService", func() {
	var topic server.MessageTopic
	var clientKey server.Key
	var subService *server.SubscriptionService
	BeforeEach(func() {
		clientKey = "some-public-key"
		authService := server.NewAuthenticationService(server.NewAuthenticationConfig())
		mockAuthService := NewAuthServiceMock(authService)
		subService = server.NewSubscriptionService(server.NewSubscriptionConfig())
		subService.AddDependency(mockAuthService)
	})
	Describe("::SubbedTopics", func() {
		When("the client is not subbed to any topics", func() {
			It("returns an empty list", func() {
				Expect(subService.SubbedTopics(clientKey).Size()).To(BeZero())
			})
		})
		When("the client is subbed to a topic", func() {
			BeforeEach(func() {
				Expect(subService.SubClient(clientKey, topic)).ToNot(HaveOccurred())
			})
			It("returns a list with that topic", func() {
				Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(1))
			})
		})
	})
	Describe("::ClientKeysSubbedToTopic", func() {
		When("no clients are subbed to the topic", func() {
			It("returns an empty list", func() {
				subbedClientKeys := subService.ClientKeysSubbedToTopic(topic)
				Expect(subbedClientKeys.Size()).To(BeZero())
			})
		})
		When("a client is subbed to the topic", func() {
			BeforeEach(func() {
				Expect(subService.SubClient(clientKey, topic)).ToNot(HaveOccurred())
			})
			It("returns a list with the subbed client key", func() {
				subbedClientKeys := subService.ClientKeysSubbedToTopic(topic)
				Expect(subbedClientKeys.Size()).To(Equal(1))
				Expect(subbedClientKeys.Has(clientKey)).To(BeTrue())
			})
		})
	})
	Describe("::SubClient", func() {
		When("the client is already subscribed", func() {
			BeforeEach(func() {
				Expect(subService.SubClient(clientKey, topic)).ToNot(HaveOccurred())
			})
			It("returns an error", func() {
				Expect(subService.SubClient(clientKey, topic)).To(HaveOccurred())
			})
		})
		When("the client is not subscribed", func() {
			It("subscribes the client", func() {
				Expect(subService.SubClient(clientKey, topic)).ToNot(HaveOccurred())
				Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(1))
				Expect(subService.ClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
			})
		})
	})
	Describe("::UnsubClient", func() {
		When("the client is not subscribed", func() {
			It("returns an error", func() {
				Expect(subService.UnsubClient(clientKey, topic)).To(HaveOccurred())
			})
		})
		When("the client is subscribed", func() {
			BeforeEach(func() {
				Expect(subService.SubClient(clientKey, topic)).ToNot(HaveOccurred())
			})
			It("unsubscribes the client", func() {
				Expect(subService.UnsubClient(clientKey, topic)).ToNot(HaveOccurred())
				Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(0))
				Expect(subService.ClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
			})
		})
	})
	Describe("::UnsubClientFromAll", func() {
		BeforeEach(func() {
			Expect(subService.SubClient(clientKey, topic)).ToNot(HaveOccurred())
			Expect(subService.SubClient(clientKey, "other-topic")).ToNot(HaveOccurred())
			Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(2))
		})
		It("unsubs client from all topics", func() {
			subService.UnsubClientFromAll(clientKey)
			Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(0))
		})
	})
	//Describe("::SubClient", func() {
	//	var topic MessageTopic
	//	BeforeEach(func() {
	//		topic = "auth"
	//		subbedTopics := subService.SubbedTopics(clientKey)
	//		Expect(subbedTopics.Size()).To(Equal(0))
	//		subbedClientKeys := subService.ClientKeysSubbedToTopic(topic)
	//		Expect(subbedClientKeys.Size()).To(Equal(0))
	//	})
	//	It("adds the topic to the set of subscribed topics keyed by the client key", func() {
	//		err := subService.SubClient(clientKey, topic)
	//		Expect(err).ToNot(HaveOccurred())
	//		Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(1))
	//		Expect(subService.ClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
	//	})
	//	When("the client is already subscribed", func() {
	//		BeforeEach(func() {
	//			err := subService.SubClient(clientKey, topic)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(1))
	//			Expect(subService.ClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
	//		})
	//		It("returns an error", func() {
	//			err := subService.SubClient(clientKey, topic)
	//			Expect(err).To(HaveOccurred())
	//			Expect(err).To(Equal(fmt.Errorf("client %s already subscribed to topic %s", clientKey, topic)))
	//		})
	//	})
	//})
	//Describe("::UnsubClient", func() {
	//	var topic MessageTopic
	//	BeforeEach(func() {
	//		topic = "auth"
	//		err := subService.SubClient(clientKey, topic)
	//		Expect(err).ToNot(HaveOccurred())
	//		Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(1))
	//		Expect(subService.ClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
	//	})
	//	It("removes the topic from the clients subscribed topic", func() {
	//		err := subService.UnsubClient(clientKey, topic)
	//		Expect(err).ToNot(HaveOccurred())
	//		Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(0))
	//		Expect(subService.ClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
	//	})
	//	When("the client is not subscribed to the topic", func() {
	//		BeforeEach(func() {
	//			err := subService.UnsubClient(clientKey, topic)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(0))
	//			Expect(subService.ClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
	//		})
	//		It("returns an error", func() {
	//			err := subService.UnsubClient(clientKey, topic)
	//			Expect(err).To(HaveOccurred())
	//			Expect(err).To(Equal(fmt.Errorf("client %s not subscribed to topic %s", clientKey, topic)))
	//		})
	//	})
	//})
	//Describe("::UnsubClientFromAll", func() {
	//	var topicA MessageTopic
	//	var topicB MessageTopic
	//	BeforeEach(func() {
	//		topicA = "auth"
	//		topicB = "findMatch"
	//		err := subService.SubClient(clientKey, topicA)
	//		Expect(err).ToNot(HaveOccurred())
	//		Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(1))
	//		Expect(subService.ClientKeysSubbedToTopic(topicA).Size()).To(Equal(1))
	//		err = subService.SubClient(clientKey, topicB)
	//		Expect(err).ToNot(HaveOccurred())
	//		Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(2))
	//		Expect(subService.ClientKeysSubbedToTopic(topicB).Size()).To(Equal(1))
	//	})
	//	It("removes both topics from the client's subscriptions", func() {
	//		subService.UnsubClientFromAll(clientKey)
	//		Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(0))
	//		Expect(subService.ClientKeysSubbedToTopic(topicA).Size()).To(Equal(0))
	//		Expect(subService.ClientKeysSubbedToTopic(topicB).Size()).To(Equal(0))
	//	})
	//	When("the client has never subscribed to any topics", func() {
	//		It("does not remove subscribed topics from other clients", func() {
	//			subService.UnsubClientFromAll("some-random-client-key")
	//			Expect(subService.SubbedTopics(clientKey).Size()).To(Equal(2))
	//		})
	//	})
	//})
	//Describe("::ClientKeysSubbedToTopic", func() {
	//	When("the topic has never been subscribed to", func() {
	//		It("initializes an empty set and returns it", func() {
	//			subbedClientKeys := subService.ClientKeysSubbedToTopic("auth")
	//			Expect(*subbedClientKeys).To(BeAssignableToTypeOf(set.Set[string]{}))
	//			Expect(subbedClientKeys.Size()).To(Equal(0))
	//		})
	//	})
	//	When("the topic has been subscribed to", func() {
	//		BeforeEach(func() {
	//			subErr := subService.SubClient(clientKey, "topicA")
	//			Expect(subErr).ToNot(HaveOccurred())
	//			subErr = subService.SubClient(clientKey, "topicB")
	//			Expect(subErr).ToNot(HaveOccurred())
	//		})
	//		It("returns the set of client keys subscribed to the topic", func() {
	//			subbedClientKeys := subService.ClientKeysSubbedToTopic("topicA")
	//			Expect(subbedClientKeys.Size()).To(Equal(1))
	//			Expect(subbedClientKeys.Has(clientKey)).To(BeTrue())
	//		})
	//	})
	//})
	//Describe("::SubbedTopics", func() {
	//	When("the client never subscribed to any topics", func() {
	//		It("initializes an empty set and returns it", func() {
	//			subbedTopics := subService.SubbedTopics(clientKey)
	//			Expect(subbedTopics).To(BeAssignableToTypeOf(set.Set[MessageTopic]{}))
	//			Expect(subbedTopics.Size()).To(Equal(0))
	//		})
	//	})
	//	When("the client has subscribed to topics", func() {
	//		BeforeEach(func() {
	//			subErr := subService.SubClient(clientKey, "topicA")
	//			Expect(subErr).ToNot(HaveOccurred())
	//			subErr = subService.SubClient(clientKey, "topicB")
	//			Expect(subErr).ToNot(HaveOccurred())
	//		})
	//		It("returns the set of topics the client is subscribed to", func() {
	//			subbedTopics := subService.SubbedTopics(clientKey)
	//			Expect(subbedTopics.Size()).To(Equal(2))
	//			Expect(subbedTopics.Has("topicA")).To(BeTrue())
	//			Expect(subbedTopics.Has("topicB")).To(BeTrue())
	//		})
	//	})
	//})
})
