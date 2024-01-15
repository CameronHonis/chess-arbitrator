package clients_manager_test

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/message_service"
	. "github.com/CameronHonis/chess-arbitrator/mocks"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-arbitrator/subscription_service"
	. "github.com/CameronHonis/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func BuildTestServices() *clients_manager.ClientsManager {
	loggerConfig := NewLoggerConfig()
	loggerConfig.MutedEnvs.Add(models.ENV_SERVER)
	subService := subscription_service.NewSubscriptionService(subscription_service.NewSubscriptionServiceConfig())
	mockSubService := subscription_service.NewSubServiceMock(subService)

	msgService := message_service.NewMessageHandlerService(message_service.NewMessageServiceConfig())
	mockMsgService := message_service.NewMessageServiceMock(msgService)

	authService := auth.NewAuthenticationService(auth.NewAuthServiceConfig())
	mockAuthService := NewAuthServiceMock(authService)

	loggerService := NewLoggerService(loggerConfig)
	mockLoggerService := NewLoggerServiceMock(loggerService)

	ucs := clients_manager.NewClientsManager(clients_manager.NewClientsManagerConfig())
	ucs.AddDependency(mockSubService)
	ucs.AddDependency(mockMsgService)
	ucs.AddDependency(mockAuthService)
	ucs.AddDependency(mockLoggerService)

	return ucs
}

type TestMessageContentType struct {
	SomePayload string `json:"somePayload"`
}

var _ = Describe("ClientsManager", func() {
	var uc *clients_manager.ClientsManager
	var client *models.Client
	BeforeEach(func() {
		uc = BuildTestServices()
		client = auth.CreateClient(nil, nil)
	})
	Describe("::AddClient", func() {
		When("the client hasn't been added", func() {
			It("adds the client to the map", func() {
				Expect(uc.AddClient(client)).ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).Error().ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).To(BeAssignableToTypeOf(&models.Client{}))
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
				subServiceMock := uc.SubService.(*subscription_service.SubServiceMock)
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
		var topic models.MessageTopic
		var msg *models.Message
		BeforeEach(func() {
			topic = "some-topic"
			msg = &models.Message{
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
			subServiceMock := uc.SubService.(*subscription_service.SubServiceMock)
			clientKeysSubbedToTopicCallCount := subServiceMock.MethodCallCount("ClientKeysSubbedToTopic")
			Expect(clientKeysSubbedToTopicCallCount).To(Equal(1))
			clientKeysSubbedToTopicCallArgs := subServiceMock.LastCallArgs("ClientKeysSubbedToTopic")
			Expect(clientKeysSubbedToTopicCallArgs[0]).To(Equal(topic))
		})
		// TODO: implement once stub generator exists
		//When("subscribers are listening on the topic", func() {
		//})
	})
})
