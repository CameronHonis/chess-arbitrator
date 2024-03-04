package clients_manager_test

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	cm "github.com/CameronHonis/chess-arbitrator/clients_manager"
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

func CreateServices(ctrl *gomock.Controller) *cm.ClientsManager {
	subServiceMock := mocks.NewMockSubscriptionServiceI(ctrl)
	subServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	subServiceMock.EXPECT().Build().AnyTimes()

	authServiceMock := mocks.NewMockAuthenticationServiceI(ctrl)
	authServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	authServiceMock.EXPECT().Build().AnyTimes()

	loggerServiceMock := mocks.NewMockLoggerServiceI(ctrl)
	loggerServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	loggerServiceMock.EXPECT().Build().AnyTimes()
	loggerServiceMock.EXPECT().Log(gomock.All(), gomock.Any()).AnyTimes()
	loggerServiceMock.EXPECT().LogRed(gomock.All(), gomock.Any()).AnyTimes()

	matchmakingMock := mocks.NewMockMatchmakingServiceI(ctrl)
	matchmakingMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	matchmakingMock.EXPECT().Build().AnyTimes()

	matcherServiceMock := mocks.NewMockMatcherServiceI(ctrl)
	matcherServiceMock.EXPECT().SetParent(gomock.All()).AnyTimes()
	matcherServiceMock.EXPECT().Build().AnyTimes()

	ucs := cm.NewClientsManager(cm.NewClientsManagerConfig(make(map[models.ContentType]cm.MessageHandler)))
	ucs.AddDependency(subServiceMock)
	ucs.AddDependency(authServiceMock)
	ucs.AddDependency(loggerServiceMock)
	ucs.AddDependency(matchmakingMock)
	ucs.AddDependency(matcherServiceMock)

	return ucs
}

type TestMessageContentType struct {
	SomePayload string `json:"somePayload"`
}

var _ = Describe("methods", func() {
	var subServiceMock *mocks.MockSubscriptionServiceI
	var eventCatcher *test_helpers.EventCatcher
	var uc *cm.ClientsManager
	BeforeEach(func() {
		ctrl := gomock.NewController(T, gomock.WithOverridableExpectations())
		uc = CreateServices(ctrl)
		eventCatcher = test_helpers.NewEventCatcher()
		eventCatcher.AddDependency(uc)
		subServiceMock = uc.SubService.(*mocks.MockSubscriptionServiceI)
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
		})
		// TODO: implement once stub generator exists
		//When("subscribers are listening on the topic", func() {
		//})
	})
})

var _ = Describe("side effects", func() {
	var uc *cm.ClientsManager
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
			cm.OnClientCreated = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnClientCreated", func() {
			ev := cm.NewClientCreatedEvent(models.NewAuthCreds("client1", "privateKey1", models.PLEB))
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
	When("a ROLE_SWITCH_GRANTED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			cm.OnUpgradeAuthGranted = func(_ service.ServiceI, _ service.EventI) bool {
				eventHandlerCalled = true
				return true
			}
		})
		It("calls OnUpgradeAuthGranted", func() {
			ev := auth.NewRoleSwitchGrantedEvent("client", "role")
			uc.Dispatch(ev)
			Eventually(func() bool {
				return eventHandlerCalled
			}).Should(BeTrue())
		})
	})
	When("a CHALLENGE_REQUEST_FAILED event is dispatched", func() {
		var eventHandlerCalled bool
		BeforeEach(func() {
			cm.OnChallengeRequestFailed = func(_ service.ServiceI, _ service.EventI) bool {
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
			cm.OnChallengeCreated = func(_ service.ServiceI, _ service.EventI) bool {
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
			cm.OnChallengeRevoked = func(_ service.ServiceI, _ service.EventI) bool {
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
			cm.OnChallengeDenied = func(_ service.ServiceI, _ service.EventI) bool {
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
