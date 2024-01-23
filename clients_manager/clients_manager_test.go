package clients_manager_test

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/clients_manager"
	"github.com/CameronHonis/chess-arbitrator/helpers/mocks"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service"
	"github.com/CameronHonis/service/test_helpers"
	"github.com/CameronHonis/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func CreateServices(ctrl *gomock.Controller) *clients_manager.ClientsManager {
	subServiceMock := mocks.NewMockSubscriptionServiceI(ctrl)
	subServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	subServiceMock.EXPECT().Build().AnyTimes()

	msgServiceMock := mocks.NewMockMessageServiceI(ctrl)
	msgServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	msgServiceMock.EXPECT().Build().AnyTimes()

	authServiceMock := mocks.NewMockAuthenticationServiceI(ctrl)
	authServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	authServiceMock.EXPECT().Build().AnyTimes()

	loggerServiceMock := mocks.NewMockLoggerServiceI(ctrl)
	loggerServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	loggerServiceMock.EXPECT().Build().AnyTimes()
	loggerServiceMock.EXPECT().Log(gomock.All(), gomock.Any()).AnyTimes()
	loggerServiceMock.EXPECT().LogRed(gomock.All(), gomock.Any()).AnyTimes()

	ucs := clients_manager.NewClientsManager(clients_manager.NewClientsManagerConfig(false))
	ucs.AddDependency(subServiceMock)
	ucs.AddDependency(msgServiceMock)
	ucs.AddDependency(authServiceMock)
	ucs.AddDependency(loggerServiceMock)

	return ucs
}

type TestMessageContentType struct {
	SomePayload string `json:"somePayload"`
}

var _ = Describe("methods", func() {
	var subServiceMock *mocks.MockSubscriptionServiceI
	var eventCatcher *test_helpers.EventCatcher
	var uc *clients_manager.ClientsManager
	var client *models.Client
	BeforeEach(func() {
		ctrl := gomock.NewController(T, gomock.WithOverridableExpectations())
		uc = CreateServices(ctrl)
		eventCatcher = test_helpers.NewEventCatcher()
		eventCatcher.AddDependency(uc)
		subServiceMock = uc.SubService.(*mocks.MockSubscriptionServiceI)
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

var _ = Describe("side effects", func() {
	var uc *clients_manager.ClientsManager
	BeforeEach(func() {
		ctrl := gomock.NewController(T, gomock.WithOverridableExpectations())
		uc = CreateServices(ctrl)
	})
	JustBeforeEach(func() {
		uc.Build()
	})
	When("a CLIENT_CREATED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			clients_manager.OnClientCreated = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnClientCreated", func() {
			ev := clients_manager.NewClientCreatedEvent(models.NewClient("client1", "privateKey1", nil, nil))
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
	When("a AUTH_UPGRADE_GRANTED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			clients_manager.OnUpgradeAuthGranted = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnUpgradeAuthGranted", func() {
			ev := auth.NewAuthUpgradeGrantedEvent("client", "role")
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
	When("a CHALLENGE_REQUEST_FAILED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			clients_manager.OnChallengeRequestFailed = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnChallengeRequestFailed", func() {
			ev := matcher.NewChallengeRequestFailedEvent(&models.Challenge{}, "")
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
	When("a CHALLENGE_CREATED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			clients_manager.OnChallengeCreated = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnChallengeRequestFailed", func() {
			ev := matcher.NewChallengeCreatedEvent(&models.Challenge{})
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
	When("a CHALLENGE_REVOKED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			clients_manager.OnChallengeRevoked = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnChallengeRevoked", func() {
			ev := matcher.NewChallengeRevokedEvent(&models.Challenge{})
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
	When("a CHALLENGE_DENIED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			clients_manager.OnChallengeDenied = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnChallengeDenied", func() {
			ev := matcher.NewChallengeDeniedEvent(&models.Challenge{})
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
})
