package matcher_test

import (
	"fmt"
	mock_auth "github.com/CameronHonis/chess-arbitrator/auth/mock"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	mock_log "github.com/CameronHonis/log/mock"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

type EventCatcher struct {
	Service
	__Dependencies__ Marker
	ListeningTo      ServiceI

	__State__    Marker
	evs          []EventI
	evsByVariant map[EventVariant][]EventI
}

func NewEventCatcher() *EventCatcher {
	ec := &EventCatcher{
		evs:          make([]EventI, 0),
		evsByVariant: make(map[EventVariant][]EventI),
	}
	ec.Service = *NewService(ec, nil)
	catcher := func(ev EventI) bool {
		ec.CatchEvent(ev)
		return false
	}
	ec.AddEventListener(ALL_EVENTS, catcher)
	return ec
}

func (ec *EventCatcher) CatchEvent(ev EventI) {
	ec.evs = append(ec.evs, ev)
	if _, ok := ec.evsByVariant[ev.Variant()]; !ok {
		ec.evsByVariant[ev.Variant()] = make([]EventI, 0)
	}
	ec.evsByVariant[ev.Variant()] = append(ec.evsByVariant[ev.Variant()], ev)
}

func (ec *EventCatcher) LastEvent() EventI {
	if len(ec.evs) == 0 {
		panic("no events have been caught")
	}
	return ec.evs[len(ec.evs)-1]
}

func (ec *EventCatcher) LastEventByVariant(eVar EventVariant) EventI {
	evs, ok := ec.evsByVariant[eVar]
	if !ok {
		panic(fmt.Sprintf("no events with variant %s have been caught", eVar))
	}
	return evs[len(evs)-1]
}

func (ec *EventCatcher) EventsCount() int {
	return len(ec.evs)
}

func (ec *EventCatcher) EventsByVariantCount(eVar EventVariant) int {
	evs, ok := ec.evsByVariant[eVar]
	if !ok {
		return 0
	}
	return len(evs)
}

func (ec *EventCatcher) NthEvent(idx int) EventI {
	if idx >= len(ec.evs) {
		panic(fmt.Sprintf("idx %d exceeds bounds of caught events (size %d)", idx, len(ec.evs)))
	}
	return ec.evs[idx]
}

func (ec *EventCatcher) NthEventByVariant(eVar EventVariant, idx int) EventI {
	evs, ok := ec.evsByVariant[eVar]
	if !ok {
		panic(fmt.Sprintf("no %s events have been caught", eVar))
	}
	if idx >= len(evs) {
		panic(fmt.Sprintf("idx %d exceeds bounds of caught %s events (size %d)", idx, eVar, len(evs)))
	}
	return evs[idx]
}

func BuildServices(ctrl *gomock.Controller) *matcher.MatcherService {
	authServiceMock := mock_auth.NewMockAuthenticationServiceI(ctrl)
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

	logServiceMock := mock_log.NewMockLoggerServiceI(ctrl)
	logServiceMock.EXPECT().SetParent(gomock.Any()).AnyTimes()
	logServiceMock.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
	logServiceMock.EXPECT().LogRed(gomock.Any(), gomock.Any()).AnyTimes()

	matcher := matcher.NewMatcherService(matcher.NewMatcherServiceConfig())
	matcher.AddDependency(authServiceMock)
	matcher.AddDependency(logServiceMock)
	return matcher
}

var _ = Describe("MatcherService", func() {
	var matcherService *matcher.MatcherService
	var eventCatcher *EventCatcher
	BeforeEach(func() {
		ctrl := gomock.NewController(T, gomock.WithOverridableExpectations())
		matcherService = BuildServices(ctrl)
		eventCatcher = NewEventCatcher()
		eventCatcher.AddDependency(matcherService)
	})
	Describe("AddMatch", func() {
		var match *models.Match
		BeforeEach(func() {
			match = models.NewMatch(
				"client1",
				"client2",
				models.NewBulletTimeControl(),
			)
		})
		Describe("when one of the players in the proposed matcherService is already in a matcherService", func() {
			BeforeEach(func() {
				ongoingMatch := models.NewMatch(
					"client1",
					"client3",
					models.NewBulletTimeControl(),
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
			match = models.NewMatch("client1", "client2", models.NewBulletTimeControl())
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
	//Describe("SetMatch", func() {
	//	var newMatch *Match
	//	BeforeEach(func() {
	//		newMatch = NewMatch("client1", "client2", NewBulletTimeControl())
	//		newMatch.WhiteClientKey = "client1"
	//		newMatch.BlackClientKey = "client2"
	//		move := chess.Move{chess.WHITE_PAWN, &chess.Square{2, 4}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
	//		newBoard := chess.GetBoardFromMove(newMatch.Board, &move)
	//		newTime := time.Now().Add(time.Second * -10)
	//		newMatch.Board = newBoard
	//		newMatch.LastMoveTime = &newTime
	//	})
	//	Describe("when the matcherService exists", func() {
	//		var prevMatch *Match
	//		BeforeEach(func() {
	//			prevMatch = NewMatch("client1", "client2", NewBulletTimeControl())
	//			prevMatch.WhiteClientKey = "client1"
	//			prevMatch.BlackClientKey = "client2"
	//			prevMatch.Uuid = newMatch.Uuid
	//			matcherService.matchByMatchId[prevMatch.Uuid] = prevMatch
	//			matcherService.matchIdByClientKey[prevMatch.WhiteClientKey] = prevMatch.Uuid
	//			matcherService.matchIdByClientKey[prevMatch.BlackClientKey] = prevMatch.Uuid
	//		})
	//		It("updates the matcherService", func() {
	//			err := matcherService.setMatch(newMatch)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(matcherService.matchByMatchId[newMatch.Uuid]).To(Equal(newMatch))
	//		})
	//		It("emits a matcherService updated event", func() {
	//			err := matcherService.setMatch(newMatch)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(el.Events).To(HaveLen(1))
	//			Expect(el.Events[0].Variant()).To(Equal(MATCH_UPDATED))
	//			Expect(el.Events[0].Payload()).To(Equal(newMatch))
	//		})
	//		Describe("when the new matcherService differs by client id", func() {
	//			BeforeEach(func() {
	//				newMatch = NewMatch("other-client1", "client2", NewBulletTimeControl())
	//				newMatch.WhiteClientKey = "other-client1"
	//				newMatch.BlackClientKey = "client2"
	//				newMatch.Uuid = prevMatch.Uuid
	//			})
	//			It("returns an error", func() {
	//				err := matcherService.setMatch(newMatch)
	//				Expect(err).To(HaveOccurred())
	//			})
	//		})
	//		Describe("when the new matcherService differs by time control", func() {
	//			BeforeEach(func() {
	//				newMatch = NewMatch("client1", "client2", NewRapidTimeControl())
	//				newMatch.WhiteClientKey = "client1"
	//				newMatch.BlackClientKey = "client2"
	//				newMatch.Uuid = prevMatch.Uuid
	//			})
	//			It("returns an error", func() {
	//				err := matcherService.setMatch(newMatch)
	//				Expect(err).To(HaveOccurred())
	//			})
	//		})
	//	})
	//	Describe("when the matcherService does not exist", func() {
	//		BeforeEach(func() {
	//			Expect(matcherService.matchByMatchId).ToNot(HaveKey(newMatch.Uuid))
	//		})
	//		It("returns an error", func() {
	//			err := matcherService.setMatch(newMatch)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//})
	//
	//Describe("ChallengeClient", func() {
	//
	//})
	//Describe("ExecuteMove", func() {
	//	var matcherService *Match
	//	var move chess.Move
	//	BeforeEach(func() {
	//		matcherService = NewMatch("client1", "client2", NewBulletTimeControl())
	//		move = chess.Move{
	//			Piece:               chess.WHITE_PAWN,
	//			StartSquare:         &chess.Square{Rank: 2, File: 4},
	//			EndSquare:           &chess.Square{Rank: 4, File: 4},
	//			CapturedPiece:       chess.EMPTY,
	//			KingCheckingSquares: make([]*chess.Square, 0),
	//			PawnUpgradedTo:      chess.EMPTY,
	//		}
	//	})
	//	Describe("when the matcherService doesnt exist", func() {
	//		BeforeEach(func() {
	//			_, getMatchErr := matcherService.MatchById(matcherService.Uuid)
	//			Expect(getMatchErr).To(HaveOccurred())
	//		})
	//		It("returns an error", func() {
	//			err := matcherService.ExecuteMove(matcherService.Uuid, &move)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//	Describe("when the matcherService exists", func() {
	//		BeforeEach(func() {
	//			matcherService = NewMatch("client1", "client2", NewBulletTimeControl())
	//			addMatchErr := matcherService.AddMatch(matcherService)
	//			Expect(addMatchErr).ToNot(HaveOccurred())
	//			Expect(matcherService.matchByMatchId).To(HaveKey(matcherService.Uuid))
	//		})
	//		Describe("when the move is illegal", func() {
	//			BeforeEach(func() {
	//				move = chess.Move{chess.WHITE_PAWN, &chess.Square{8, 8}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
	//			})
	//			It("returns an error", func() {
	//				err := matcherService.ExecuteMove(matcherService.Uuid, &move)
	//				Expect(err).To(HaveOccurred())
	//			})
	//		})
	//		It("Updates the matcherService", func() {
	//			err := matcherService.ExecuteMove(matcherService.Uuid, &move)
	//			Expect(err).ToNot(HaveOccurred())
	//			expBoard := chess.GetBoardFromMove(matcherService.Board, &move)
	//			Expect(matcherService.matchByMatchId[matcherService.Uuid].Board).To(Equal(expBoard))
	//		})
	//	})
	//})
})
