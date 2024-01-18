package clients_manager_test

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	mock_auth "github.com/CameronHonis/chess-arbitrator/auth/mock"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/models"
	mock_msg_service "github.com/CameronHonis/chess-arbitrator/msg_service/mock"
	mock_sub_service "github.com/CameronHonis/chess-arbitrator/sub_service/mock"
	mock_log "github.com/CameronHonis/log/mock"
	"github.com/CameronHonis/service/test_helpers"
	"github.com/CameronHonis/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func BuildTestServices(ctrl *gomock.Controller) *clients_manager.ClientsManager {
	subServiceMock := mock_sub_service.NewMockSubscriptionServiceI(ctrl)
	subServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()

	msgServiceMock := mock_msg_service.NewMockMessageServiceI(ctrl)
	msgServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()

	authServiceMock := mock_auth.NewMockAuthenticationServiceI(ctrl)
	authServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()

	loggerServiceMock := mock_log.NewMockLoggerServiceI(ctrl)
	loggerServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	loggerServiceMock.EXPECT().Log(gomock.All(), gomock.Any()).AnyTimes()
	loggerServiceMock.EXPECT().LogRed(gomock.All(), gomock.Any()).AnyTimes()

	ucs := clients_manager.NewClientsManager(clients_manager.NewClientsManagerConfig())
	ucs.AddDependency(subServiceMock)
	ucs.AddDependency(msgServiceMock)
	ucs.AddDependency(authServiceMock)
	ucs.AddDependency(loggerServiceMock)

	return ucs
}

type TestMessageContentType struct {
	SomePayload string `json:"somePayload"`
}

var _ = Describe("ClientsManager", func() {
	var subServiceMock *mock_sub_service.MockSubscriptionServiceI
	var eventCatcher *test_helpers.EventCatcher
	var uc *clients_manager.ClientsManager
	var client *models.Client
	BeforeEach(func() {
		ctrl := gomock.NewController(T, gomock.WithOverridableExpectations())
		uc = BuildTestServices(ctrl)
		eventCatcher = test_helpers.NewEventCatcher()
		eventCatcher.AddDependency(uc)
		subServiceMock = uc.SubService.(*mock_sub_service.MockSubscriptionServiceI)
		client = auth.CreateClient(nil, nil)
	})
	Describe("::AddClient", func() {
		BeforeEach(func() {
			eventCount := eventCatcher.EventsByVariantCount(clients_manager.CLIENT_CREATED)
			Expect(eventCount).To(Equal(0))
		})
		When("the client hasn't been added", func() {
			It("adds the client to the map", func() {
				Expect(uc.AddClient(client)).ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).Error().ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).To(BeAssignableToTypeOf(&models.Client{}))
			})
			It("emits a CLIENT_CREATED event", func() {
				Expect(uc.AddClient(client)).ToNot(HaveOccurred())
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(clients_manager.CLIENT_CREATED)
				}).Should(Equal(1))
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
				subServiceMock.EXPECT().UnsubClientFromAll(gomock.All()).AnyTimes()
			})
			It("removes the client from the client map", func() {
				Expect(uc.RemoveClient(client)).ToNot(HaveOccurred())
				Expect(uc.GetClient(client.PublicKey())).Error().To(HaveOccurred())
			})
			It("unsubs client from all topics", func() {
				subServiceMock.EXPECT().UnsubClientFromAll(client.PublicKey())
				Expect(uc.RemoveClient(client)).ToNot(HaveOccurred())
				//unsubAllCount := subServiceMock.MethodCallCount("UnsubClientFromAll")
				//Expect(unsubAllCount).To(Equal(1))
				//unsubAllArgs := subServiceMock.LastCallArgs("UnsubClientFromAll")
				//Expect(unsubAllArgs[0]).To(Equal(client.PublicKey()))
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
			subServiceMock.EXPECT().ClientKeysSubbedToTopic(gomock.Eq(topic)).Return(set.EmptySet[models.Key]())
			uc.BroadcastMessage(msg)
			//subServiceMock := uc.SubService.(*sub_service.SubServiceMock)
			//clientKeysSubbedToTopicCallCount := subServiceMock.MethodCallCount("ClientKeysSubbedToTopic")
			//Expect(clientKeysSubbedToTopicCallCount).To(Equal(1))
			//clientKeysSubbedToTopicCallArgs := subServiceMock.LastCallArgs("ClientKeysSubbedToTopic")
			//Expect(clientKeysSubbedToTopicCallArgs[0]).To(Equal(topic))
		})
		// TODO: implement once stub generator exists
		//When("subscribers are listening on the topic", func() {
		//})
	})
})
