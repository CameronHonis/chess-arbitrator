package matcher_test

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/auth_service"
	"github.com/CameronHonis/chess-arbitrator/matcher"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	. "github.com/CameronHonis/service"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

var _ = Describe("MatcherService", func() {
	var m *matcher.MatcherService
	var eventCatcher *EventCatcher
	BeforeEach(func() {
		realLoggerService := log.NewLoggerService(log.NewLoggerConfig())
		realAuthService := auth_service.NewAuthenticationService(nil)
		mockAuthService := auth_service.NewAuthServiceMock(realAuthService)
		getRoleStub := func(rec *auth_service.AuthenticationService, clientKey models.Key) (models.RoleName, error) {
			roleName := map[models.Key]models.RoleName{
				"client1": models.PLEB,
				"client2": models.PLEB,
				"client3": models.PLEB,
			}[clientKey]
			if roleName == "" {
				return "", fmt.Errorf("no role for %s", clientKey)
			}
			return roleName, nil
		}
		mockAuthService.Stub("GetRole", getRoleStub)
		m = matcher.NewMatcherService(nil)
		m.AddDependency(mockAuthService)
		m.AddDependency(realLoggerService)
		eventCatcher = NewEventCatcher()
		eventCatcher.AddDependency(m)
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
		Describe("when one of the players in the proposed matcher is already in a matcher", func() {
			BeforeEach(func() {
				ongoingMatch := models.NewMatch(
					"client1",
					"client3",
					models.NewBulletTimeControl(),
				)
				Expect(m.AddMatch(ongoingMatch)).To(Succeed())
			})
			It("returns an error", func() {
				Expect(m.AddMatch(match)).ToNot(Succeed())
			})
		})
		It("adds the matcher to the active matches", func() {
			Expect(m.AddMatch(match)).To(Succeed())
			Expect(m.MatchById(match.Uuid)).To(Equal(match))
			Expect(m.MatchByClientKey(match.BlackClientKey)).To(Equal(match))
			Expect(m.MatchByClientKey(match.WhiteClientKey)).To(Equal(match))
		})
		It("emits a matcher created event", func() {
			Expect(m.AddMatch(match)).To(Succeed())

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
		Describe("when the matcher exists", func() {
			BeforeEach(func() {
				Expect(m.AddMatch(match)).To(Succeed())
			})
			It("removes the matcher from the active matches", func() {
				Expect(m.RemoveMatch(match)).To(Succeed())
				Expect(m.MatchById(match.Uuid)).Error().To(HaveOccurred())
				Expect(m.MatchByClientKey(match.WhiteClientKey)).Error().To(HaveOccurred())
				Expect(m.MatchByClientKey(match.BlackClientKey)).Error().To(HaveOccurred())
			})
			It("emits a matcher ended event", func() {
				Expect(m.RemoveMatch(match)).To(Succeed())
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(matcher.MATCH_ENDED)
				}).Should(Equal(1))
				expEvent := matcher.NewMatchEndedEvent(match)
				actualEvent := eventCatcher.LastEventByVariant(matcher.MATCH_ENDED)
				Expect(actualEvent).To(BeEquivalentTo(expEvent))
			})
		})
		Describe("when the matcher does not exist", func() {
			It("returns an error", func() {
				Expect(m.RemoveMatch(match)).To(HaveOccurred())
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
	//	Describe("when the matcher exists", func() {
	//		var prevMatch *Match
	//		BeforeEach(func() {
	//			prevMatch = NewMatch("client1", "client2", NewBulletTimeControl())
	//			prevMatch.WhiteClientKey = "client1"
	//			prevMatch.BlackClientKey = "client2"
	//			prevMatch.Uuid = newMatch.Uuid
	//			m.matchByMatchId[prevMatch.Uuid] = prevMatch
	//			m.matchIdByClientKey[prevMatch.WhiteClientKey] = prevMatch.Uuid
	//			m.matchIdByClientKey[prevMatch.BlackClientKey] = prevMatch.Uuid
	//		})
	//		It("updates the matcher", func() {
	//			err := m.setMatch(newMatch)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(m.matchByMatchId[newMatch.Uuid]).To(Equal(newMatch))
	//		})
	//		It("emits a matcher updated event", func() {
	//			err := m.setMatch(newMatch)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(el.Events).To(HaveLen(1))
	//			Expect(el.Events[0].Variant()).To(Equal(MATCH_UPDATED))
	//			Expect(el.Events[0].Payload()).To(Equal(newMatch))
	//		})
	//		Describe("when the new matcher differs by client id", func() {
	//			BeforeEach(func() {
	//				newMatch = NewMatch("other-client1", "client2", NewBulletTimeControl())
	//				newMatch.WhiteClientKey = "other-client1"
	//				newMatch.BlackClientKey = "client2"
	//				newMatch.Uuid = prevMatch.Uuid
	//			})
	//			It("returns an error", func() {
	//				err := m.setMatch(newMatch)
	//				Expect(err).To(HaveOccurred())
	//			})
	//		})
	//		Describe("when the new matcher differs by time control", func() {
	//			BeforeEach(func() {
	//				newMatch = NewMatch("client1", "client2", NewRapidTimeControl())
	//				newMatch.WhiteClientKey = "client1"
	//				newMatch.BlackClientKey = "client2"
	//				newMatch.Uuid = prevMatch.Uuid
	//			})
	//			It("returns an error", func() {
	//				err := m.setMatch(newMatch)
	//				Expect(err).To(HaveOccurred())
	//			})
	//		})
	//	})
	//	Describe("when the matcher does not exist", func() {
	//		BeforeEach(func() {
	//			Expect(m.matchByMatchId).ToNot(HaveKey(newMatch.Uuid))
	//		})
	//		It("returns an error", func() {
	//			err := m.setMatch(newMatch)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//})
	//
	//Describe("ChallengeClient", func() {
	//
	//})
	//Describe("ExecuteMove", func() {
	//	var matcher *Match
	//	var move chess.Move
	//	BeforeEach(func() {
	//		matcher = NewMatch("client1", "client2", NewBulletTimeControl())
	//		move = chess.Move{
	//			Piece:               chess.WHITE_PAWN,
	//			StartSquare:         &chess.Square{Rank: 2, File: 4},
	//			EndSquare:           &chess.Square{Rank: 4, File: 4},
	//			CapturedPiece:       chess.EMPTY,
	//			KingCheckingSquares: make([]*chess.Square, 0),
	//			PawnUpgradedTo:      chess.EMPTY,
	//		}
	//	})
	//	Describe("when the matcher doesnt exist", func() {
	//		BeforeEach(func() {
	//			_, getMatchErr := m.MatchById(matcher.Uuid)
	//			Expect(getMatchErr).To(HaveOccurred())
	//		})
	//		It("returns an error", func() {
	//			err := m.ExecuteMove(matcher.Uuid, &move)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//	Describe("when the matcher exists", func() {
	//		BeforeEach(func() {
	//			matcher = NewMatch("client1", "client2", NewBulletTimeControl())
	//			addMatchErr := m.AddMatch(matcher)
	//			Expect(addMatchErr).ToNot(HaveOccurred())
	//			Expect(m.matchByMatchId).To(HaveKey(matcher.Uuid))
	//		})
	//		Describe("when the move is illegal", func() {
	//			BeforeEach(func() {
	//				move = chess.Move{chess.WHITE_PAWN, &chess.Square{8, 8}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
	//			})
	//			It("returns an error", func() {
	//				err := m.ExecuteMove(matcher.Uuid, &move)
	//				Expect(err).To(HaveOccurred())
	//			})
	//		})
	//		It("Updates the matcher", func() {
	//			err := m.ExecuteMove(matcher.Uuid, &move)
	//			Expect(err).ToNot(HaveOccurred())
	//			expBoard := chess.GetBoardFromMove(matcher.Board, &move)
	//			Expect(m.matchByMatchId[matcher.Uuid].Board).To(Equal(expBoard))
	//		})
	//	})
	//})
})
