package server_test

import (
	"github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/CameronHonis/chess-arbitrator/server/mocks"
	"github.com/CameronHonis/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func BuildTestServices() *server.UserClientsService {
	subService := server.NewSubscriptionService(server.NewSubscriptionConfig())
	mockSubService := NewSubServiceMock(subService)

	msgService := server.NewMessageHandlerService(server.NewMessageHandlerConfig())
	mockMsgService := NewMessageServiceMock(msgService)

	authService := server.NewAuthenticationService(server.NewAuthenticationConfig())
	mockAuthService := NewAuthServiceMock(authService)

	loggerService := log.NewLoggerService(log.NewLoggerConfig())

	ucs := server.NewUserClientsService(server.NewUserClientsConfig())
	ucs.AddDependency(mockSubService)
	ucs.AddDependency(mockMsgService)
	ucs.AddDependency(mockAuthService)
	ucs.AddDependency(loggerService)
	return ucs
}

type TestMessageContentType struct {
	SomePayload string `json:"somePayload"`
}

var _ = Describe("UserClientsService", func() {
	var uc *server.UserClientsService
	var client *server.Client
	BeforeEach(func() {
		uc = BuildTestServices()
		client = server.NewClient(nil, nil)
	})
	Describe("::AddClient", func() {
		When("the client hasn't been added", func() {
			It("adds the client to the map", func() {
				Expect(uc.AddClient(client)).ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).Error().ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).To(BeAssignableToTypeOf(&server.Client{}))
			})
		})
		When("the client already exists", func() {
			BeforeEach(func() {
				Expect(uc.AddClient(client)).ToNot(HaveOccurred())
			})
			It("returns an error", func() {
				Expect(uc.AddClient(client)).Error().To(HaveOccurred())
			})
		})
	})
	Describe("::RemoveClient", func() {
		When("a player was subscribed to a topic", func() {
			BeforeEach(func() {
				Expect(uc.AddClient(client)).ToNot(HaveOccurred())
			})
			It("removes the client from the client map", func() {
				Expect(uc.RemoveClient(client)).ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).Error().To(HaveOccurred())
			})
			It("unsubs client from all topics", func() {
				Expect(uc.RemoveClient(client)).ToNot(HaveOccurred())
				subServiceMock := uc.SubscriptionService.(*SubServiceMock)
				unsubAllCount := subServiceMock.MethodCallCount("UnsubClientFromAll")
				Expect(unsubAllCount).To(Equal(1))
				unsubAllArgs := subServiceMock.LastCallArgs("UnsubClientFromAll")
				Expect(unsubAllArgs[0]).To(Equal(client.PublicKey()))
			})
		})
		When("the player was never added", func() {
			It("returns an error", func() {
				Expect(uc.RemoveClient(client)).To(HaveOccurred())
			})
		})
	})
	Describe("::BroadcastMessage", func() {
		var topic server.MessageTopic
		var msg *server.Message
		BeforeEach(func() {
			topic = "some-topic"
			msg = &server.Message{
				SenderKey:   "some-sender-key",
				PrivateKey:  "some-private-key",
				Topic:       topic,
				ContentType: "TEST_MESSAGE",
				Content: &TestMessageContentType{
					SomePayload: "some-payload-text",
				},
			}
		})
		It("queries the subscribers on the message topic", func() {
			uc.BroadcastMessage(msg)
			subServiceMock := uc.SubscriptionService.(*SubServiceMock)
			clientKeysSubbedToTopicCallCount := subServiceMock.MethodCallCount("ClientKeysSubbedToTopic")
			Expect(clientKeysSubbedToTopicCallCount).To(Equal(1))
			clientKeysSubbedToTopicCallArgs := subServiceMock.LastCallArgs("ClientKeysSubbedToTopic")
			Expect(clientKeysSubbedToTopicCallArgs[0]).To(Equal(topic))
		})
	})
})
