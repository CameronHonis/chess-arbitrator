package matcher_test

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/helpers/mocks"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/service/test_helpers"
	"github.com/CameronHonis/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"time"
)

func BuildServices(ctrl *gomock.Controller) *matcher.MatcherService {
	authServiceMock := mocks.NewMockAuthenticationServiceI(ctrl)
	authServiceMock.EXPECT().SetParent(gomock.Any()).AnyTimes()
	getRole := func(clientKey models.Key) (models.RoleName, error) {
		roleByKey := map[models.Key]models.RoleName{
			"client1": models.PLEB,
			"client2": models.PLEB,
			"client3": models.PLEB,
		}
		if role, ok := roleByKey[clientKey]; ok {
			return role, nil
		}
		return "", fmt.Errorf("client with key %s is not assigned a role", clientKey)
	}
	authServiceMock.EXPECT().GetRole(gomock.Any()).DoAndReturn(getRole).AnyTimes()

	logServiceMock := mocks.NewMockLoggerServiceI(ctrl)
	logServiceMock.EXPECT().SetParent(gomock.Any()).AnyTimes()
	logServiceMock.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
	logServiceMock.EXPECT().LogRed(gomock.Any(), gomock.Any()).AnyTimes()

	matcher_service := matcher.NewMatcherService(matcher.NewMatcherServiceConfig())
	matcher_service.AddDependency(authServiceMock)
	matcher_service.AddDependency(logServiceMock)
	return matcher_service
}

var _ = Describe("MatcherService", func() {
	var matcherService *matcher.MatcherService
	var authServiceMock *mocks.MockAuthenticationServiceI
	var eventCatcher *test_helpers.EventCatcher
	BeforeEach(func() {
		ctrl := gomock.NewController(T, gomock.WithOverridableExpectations())
		matcherService = BuildServices(ctrl)
		authServiceMock = matcherService.AuthService.(*mocks.MockAuthenticationServiceI)
		eventCatcher = test_helpers.NewEventCatcher()
		eventCatcher.AddDependency(matcherService)
	})
	Describe("AddMatch", func() {
		var match *models.Match
		BeforeEach(func() {
			match = builders.NewMatch(
				"client1",
				"client2",
				builders.NewBulletTimeControl(),
			)
		})
		Describe("when one of the players in the proposed matcherService is already in a matcherService", func() {
			BeforeEach(func() {
				ongoingMatch := builders.NewMatch(
					"client1",
					"client3",
					builders.NewBulletTimeControl(),
				)
				Expect(matcherService.AddMatch(ongoingMatch)).To(Succeed())
			})
			It("returns an error", func() {
				Expect(matcherService.AddMatch(match)).ToNot(Succeed())
			})
		})
		It("adds the matcherService to the active matches", func() {
			Expect(matcherService.AddMatch(match)).To(Succeed())
			Expect(matcherService.MatchById(match.Uuid)).To(Equal(match))
			Expect(matcherService.MatchByClientKey(match.BlackClientKey)).To(Equal(match))
			Expect(matcherService.MatchByClientKey(match.WhiteClientKey)).To(Equal(match))
		})
		It("emits a matcherService created event", func() {
			Expect(matcherService.AddMatch(match)).To(Succeed())

			Eventually(func() int {
				return eventCatcher.EventsByVariantCount(matcher.MATCH_CREATED)
			}).Should(Equal(1))
			expEvent := matcher.NewMatchCreatedEvent(match)
			actualEvent := eventCatcher.LastEventByVariant(matcher.MATCH_CREATED)
			Expect(actualEvent).To(BeEquivalentTo(expEvent))
		})
	})
	Describe("RemoveMatch", func() {
		var match *models.Match
		BeforeEach(func() {
			match = builders.NewMatch("client1", "client2", builders.NewBulletTimeControl())
		})
		Describe("when the matcherService exists", func() {
			BeforeEach(func() {
				Expect(matcherService.AddMatch(match)).To(Succeed())
			})
			It("removes the matcherService from the active matches", func() {
				Expect(matcherService.RemoveMatch(match)).To(Succeed())
				Expect(matcherService.MatchById(match.Uuid)).Error().To(HaveOccurred())
				Expect(matcherService.MatchByClientKey(match.WhiteClientKey)).Error().To(HaveOccurred())
				Expect(matcherService.MatchByClientKey(match.BlackClientKey)).Error().To(HaveOccurred())
			})
			It("emits a matcherService ended event", func() {
				Expect(matcherService.RemoveMatch(match)).To(Succeed())
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(matcher.MATCH_ENDED)
				}).Should(Equal(1))
				expEvent := matcher.NewMatchEndedEvent(match)
				actualEvent := eventCatcher.LastEventByVariant(matcher.MATCH_ENDED)
				Expect(actualEvent).To(BeEquivalentTo(expEvent))
			})
		})
		Describe("when the matcherService does not exist", func() {
			It("returns an error", func() {
				Expect(matcherService.RemoveMatch(match)).To(HaveOccurred())
			})
		})
	})
	Describe("SetMatch", func() {
		var newMatch *models.Match
		BeforeEach(func() {
			newMatch = builders.NewMatch("client1", "client2", builders.NewBulletTimeControl())
			newMatch.WhiteClientKey = "client1"
			newMatch.BlackClientKey = "client2"
			move := chess.Move{chess.WHITE_PAWN, &chess.Square{2, 4}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
			newBoard := chess.GetBoardFromMove(newMatch.Board, &move)
			newTime := time.Now().Add(time.Second * -10)
			newMatch.Board = newBoard
			newMatch.LastMoveTime = &newTime
		})
		Describe("when the match exists", func() {
			var prevMatch *models.Match
			BeforeEach(func() {
				prevMatch = builders.NewMatch("client1", "client2", builders.NewBulletTimeControl())
				prevMatch.WhiteClientKey = "client1"
				prevMatch.BlackClientKey = "client2"
				prevMatch.Uuid = newMatch.Uuid
				Expect(matcherService.AddMatch(prevMatch)).ToNot(HaveOccurred())
			})
			It("updates the match", func() {
				Expect(matcherService.SetMatch(newMatch)).ToNot(HaveOccurred())
				Expect(matcherService.MatchById(newMatch.Uuid)).To(Equal(newMatch))
			})
			It("emits a match updated event", func() {
				Expect(matcherService.SetMatch(newMatch)).ToNot(HaveOccurred())
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(matcher.MATCH_UPDATED)
				}).Should(Equal(1))
				Eventually(func() *models.Match {
					lastEvent := eventCatcher.LastEventByVariant(matcher.MATCH_UPDATED)
					if lastEvent == nil {
						return nil
					} else {
						return lastEvent.Payload().(*matcher.MatchUpdatedEventPayload).Match
					}
				}).Should(Equal(newMatch))
			})
			Describe("when the new match differs by client id", func() {
				BeforeEach(func() {
					newMatch = builders.NewMatch("other-client1", "client2", builders.NewBulletTimeControl())
					newMatch.WhiteClientKey = "other-client1"
					newMatch.BlackClientKey = "client2"
					newMatch.Uuid = prevMatch.Uuid
				})
				It("returns an error", func() {
					Expect(matcherService.SetMatch(newMatch)).To(HaveOccurred())
				})
			})
			Describe("when the new matcherService differs by time control", func() {
				BeforeEach(func() {
					newMatch = builders.NewMatch("client1", "client2", builders.NewRapidTimeControl())
					newMatch.WhiteClientKey = "client1"
					newMatch.BlackClientKey = "client2"
					newMatch.Uuid = prevMatch.Uuid
				})
				It("returns an error", func() {
					Expect(matcherService.SetMatch(newMatch)).To(HaveOccurred())
				})
			})
		})
		Describe("when the match does not exist", func() {
			BeforeEach(func() {
				Expect(matcherService.MatchById(newMatch.Uuid)).Error().To(HaveOccurred())
			})
			It("returns an error", func() {
				Expect(matcherService.SetMatch(newMatch)).To(HaveOccurred())
			})
		})
	})

	Describe("ChallengeClient", func() {
		var challenge *models.Challenge
		Describe("when the challenge is directed to a player client", func() {
			BeforeEach(func() {
				challengeBuilder := builders.NewChallengeBuilder()
				challengeBuilder.WithUuid("")
				challengeBuilder.WithChallengerKey("client1")
				challengeBuilder.WithChallengedKey("client2")
				challengeBuilder.WithIsChallengerWhite(true)
				challengeBuilder.WithTimeControl(builders.NewBulletTimeControl())
				challengeBuilder.WithIsActive(true)
				challenge = challengeBuilder.Build()
			})
			It("stores the challenge", func() {
				Expect(matcherService.RequestChallenge(challenge)).ToNot(HaveOccurred())
			})
			It("emits a challenge created event", func() {
				_ = matcherService.RequestChallenge(challenge)
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_CREATED)
				}).Should(Equal(1))
			})
			Describe("when the challenger is already in a match", func() {
				BeforeEach(func() {
					existingMatch := builders.NewMatch(
						"client1", "client3", builders.NewBlitzTimeControl(),
					)
					Expect(matcherService.AddMatch(existingMatch)).ToNot(HaveOccurred())
				})
				It("returns an error", func() {
					Expect(matcherService.RequestChallenge(challenge)).To(HaveOccurred())
				})
			})
			Describe("when the challenged is already in a match", func() {
				BeforeEach(func() {
					existingMatch := builders.NewMatch(
						"client2", "client3", builders.NewBlitzTimeControl(),
					)
					Expect(matcherService.AddMatch(existingMatch)).ToNot(HaveOccurred())
				})
				It("does not return an error", func() {
					Expect(matcherService.RequestChallenge(challenge)).ToNot(HaveOccurred())
				})
				It("stores the challenge", func() {
					_ = matcherService.RequestChallenge(challenge)
					Expect(matcherService.GetChallenge("client1", "client2")).ToNot(BeNil())
				})
				It("emits a challenge created event", func() {
					_ = matcherService.RequestChallenge(challenge)
					Eventually(func() int {
						return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_CREATED)
					}).Should(Equal(1))
				})
			})
			Describe("when a challenge to that same player has already been sent", func() {
				BeforeEach(func() {
					Expect(matcherService.RequestChallenge(challenge)).ToNot(HaveOccurred())
				})
				It("returns an error", func() {
					Expect(matcherService.RequestChallenge(challenge)).To(HaveOccurred())
				})
			})
		})
		Describe("when the challenge is directed to a bot client", func() {
			BeforeEach(func() {
				challengeBuilder := builders.NewChallengeBuilder()
				challengeBuilder.WithUuid("")
				challengeBuilder.WithChallengerKey("client1")
				challengeBuilder.WithChallengedKey("")
				challengeBuilder.WithIsChallengerWhite(true)
				challengeBuilder.WithIsChallengerBlack(false)
				challengeBuilder.WithTimeControl(builders.NewBulletTimeControl())
				challengeBuilder.WithBotName("someBot")
				challengeBuilder.WithIsActive(true)
				challenge = challengeBuilder.Build()

				authServiceMock.EXPECT().BotClientExists().Return(true).AnyTimes()
				botClientKeys := set.EmptySet[models.Key]()
				botClientKeys.Add("someBotClientKey")
				authServiceMock.EXPECT().ClientKeysByRole(gomock.Eq(models.RoleName(models.BOT))).Return(botClientKeys).AnyTimes()
			})
			It("queries for the a bot client", func() {
				authServiceMock.EXPECT().BotClientExists().Return(true)
				_ = matcherService.RequestChallenge(challenge)
			})
			It("stores the challenge", func() {
				Expect(matcherService.RequestChallenge(challenge)).ToNot(HaveOccurred())
			})
			It("emits a challenge created event", func() {
				_ = matcherService.RequestChallenge(challenge)
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_CREATED)
				}).Should(Equal(1))
			})
			Describe("when the challenger is already in a match", func() {
				BeforeEach(func() {
					match := builders.NewMatch(challenge.ChallengerKey, "client3", builders.NewBlitzTimeControl())
					Expect(matcherService.AddMatch(match)).ToNot(HaveOccurred())
				})
				It("returns an error", func() {
					Expect(matcherService.RequestChallenge(challenge)).To(HaveOccurred())
				})
			})
			Describe("no bot servers are connected", func() {
				BeforeEach(func() {
					authServiceMock.EXPECT().BotClientExists().Return(false).AnyTimes()
				})
				It("returns an error", func() {
					Expect(matcherService.RequestChallenge(challenge)).To(HaveOccurred())
				})
				It("emits a challenge creation failed event", func() {
					_ = matcherService.RequestChallenge(challenge)
					Eventually(func() int {
						return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_REQUEST_FAILED)
					}).Should(Equal(1))
				})
			})
		})
	})

	Describe("ExecuteMove", func() {
		var match *models.Match
		var move chess.Move
		BeforeEach(func() {
			match = builders.NewMatch("client1", "client2", builders.NewBulletTimeControl())
			move = chess.Move{
				Piece:               chess.WHITE_PAWN,
				StartSquare:         &chess.Square{Rank: 2, File: 4},
				EndSquare:           &chess.Square{Rank: 4, File: 4},
				CapturedPiece:       chess.EMPTY,
				KingCheckingSquares: make([]*chess.Square, 0),
				PawnUpgradedTo:      chess.EMPTY,
			}
		})
		Describe("when the match doesnt exist", func() {
			BeforeEach(func() {
				_, getMatchErr := matcherService.MatchById(match.Uuid)
				Expect(getMatchErr).To(HaveOccurred())
			})
			It("returns an error", func() {
				err := matcherService.ExecuteMove(match.Uuid, &move)
				Expect(err).To(HaveOccurred())
			})
		})
		Describe("when the match exists", func() {
			BeforeEach(func() {
				match = builders.NewMatch("client1", "client2", builders.NewBulletTimeControl())
				addMatchErr := matcherService.AddMatch(match)
				Expect(addMatchErr).ToNot(HaveOccurred())
				Expect(matcherService.MatchById(match.Uuid)).To(Equal(match))
			})
			Describe("when the move is illegal", func() {
				BeforeEach(func() {
					move = chess.Move{
						chess.WHITE_PAWN,
						&chess.Square{8, 8},
						&chess.Square{4, 4},
						chess.EMPTY,
						make([]*chess.Square, 0),
						chess.EMPTY,
					}
				})
				It("returns an error", func() {
					err := matcherService.ExecuteMove(match.Uuid, &move)
					Expect(err).To(HaveOccurred())
				})
			})
			It("Updates the match", func() {
				Expect(matcherService.ExecuteMove(match.Uuid, &move)).ToNot(HaveOccurred())
				expBoard := chess.GetBoardFromMove(match.Board, &move)
				newMatch, _ := matcherService.MatchById(match.Uuid)
				Expect(newMatch.Board).To(Equal(expBoard))
			})
		})
	})
	Describe("RevokeChallenge", func() {
		var challenge *models.Challenge
		BeforeEach(func() {
			challenge = builders.NewChallenge(
				"client1",
				"client2",
				true,
				false,
				builders.NewBlitzTimeControl(),
				"",
				false,
			)
			Expect(matcherService.RequestChallenge(challenge)).ToNot(HaveOccurred())
		})
		It("removes the challenge", func() {
			Expect(matcherService.RevokeChallenge(challenge.ChallengerKey, challenge.ChallengedKey)).ToNot(HaveOccurred())
			Expect(matcherService.GetChallenge(challenge.ChallengerKey, challenge.ChallengedKey)).Error().To(HaveOccurred())
		})
		It("emits a challenge revoked event", func() {
			_ = matcherService.RevokeChallenge(challenge.ChallengerKey, challenge.ChallengedKey)
			Eventually(func() int {
				return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_REVOKED)
			}).Should(Equal(1))
		})
	})
	Describe("DeclineChallenge", func() {
		var challenge *models.Challenge
		BeforeEach(func() {
			challenge = builders.NewChallenge(
				"client1",
				"client2",
				true,
				false,
				builders.NewBlitzTimeControl(),
				"",
				false,
			)
			Expect(matcherService.RequestChallenge(challenge)).ToNot(HaveOccurred())
		})
		It("removes the challenge", func() {
			Expect(matcherService.DeclineChallenge(challenge.ChallengerKey, challenge.ChallengedKey)).ToNot(HaveOccurred())
			Expect(matcherService.GetChallenge(challenge.ChallengerKey, challenge.ChallengedKey)).Error().To(HaveOccurred())
		})
		It("emits a challenge declined event", func() {
			_ = matcherService.DeclineChallenge(challenge.ChallengerKey, challenge.ChallengedKey)
			Eventually(func() int {
				return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_DENIED)
			}).Should(Equal(1))
		})
	})
	Describe("AcceptChallenge", func() {
		When("the challenge already exists", func() {
			var challenge *models.Challenge
			BeforeEach(func() {
				challenge = builders.NewChallenge(
					"client1",
					"client2",
					true,
					false,
					builders.NewBulletTimeControl(),
					"",
					false)
				Expect(matcherService.RequestChallenge(challenge)).ToNot(HaveOccurred())
			})
			It("dispatches a challenge accepted event", func() {
				Expect(matcherService.AcceptChallenge("client1", "client2")).ToNot(HaveOccurred())
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_ACCEPTED)
				}).Should(Equal(1))
			})
		})
		When("the challenge does not exist", func() {
			It("returns an error", func() {
				Expect(matcherService.AcceptChallenge("client1", "client2")).To(HaveOccurred())
			})
			It("dispatches a challenge accept failure event", func() {
				_ = matcherService.AcceptChallenge("client1", "client2")
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(matcher.CHALLENGE_ACCEPT_FAILED)
				}).Should(Equal(1))
			})
		})
	})
})
