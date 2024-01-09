package subscription_service_test

import (
	"github.com/CameronHonis/chess-arbitrator/auth_service"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/subscription_service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionService", func() {
	var topic models.MessageTopic
	var clientKey models.Key
	var subService *subscription_service.SubscriptionService
	BeforeEach(func() {
		clientKey = "some-public-key"
		authService := auth_service.NewAuthenticationService(auth_service.NewAuthenticationConfig())
		mockAuthService := auth_service.NewAuthServiceMock(authService)
		subService = subscription_service.NewSubscriptionService(subscription_service.NewSubscriptionConfig())
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
})
