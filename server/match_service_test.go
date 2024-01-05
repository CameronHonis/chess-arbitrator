package server_test

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/server"
	. "github.com/CameronHonis/chess-arbitrator/server/mocks"
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

var _ = Describe("MatchService", func() {
	var m *server.MatchService
	var eventCatcher *EventCatcher
	BeforeEach(func() {
		realLoggerService := log.NewLoggerService(server.GetLoggerConfig())
		realAuthService := server.NewAuthenticationService(nil)
		mockAuthService := NewAuthServiceMock(realAuthService)
		getRoleStub := func(rec *server.AuthenticationService, clientKey server.Key) (server.RoleName, error) {
			roleName := map[server.Key]server.RoleName{
				"client1": server.PLEB,
				"client2": server.PLEB,
				"client3": server.PLEB,
			}[clientKey]
			if roleName == "" {
				return "", fmt.Errorf("no role for %s", clientKey)
			}
			return roleName, nil
		}
		mockAuthService.Stub("GetRole", getRoleStub)
		m = server.NewMatchService(nil)
		m.AddDependency(mockAuthService)
		m.AddDependency(realLoggerService)
		eventCatcher = NewEventCatcher()
		eventCatcher.AddDependency(m)
	})
	Describe("AddMatch", func() {
		var match *server.Match
		BeforeEach(func() {
			match = server.NewMatch(
				"client1",
				"client2",
				server.NewBulletTimeControl(),
			)
		})
		Describe("when one of the players in the proposed match is already in a match", func() {
			BeforeEach(func() {
				ongoingMatch := server.NewMatch(
					"client1",
					"client3",
					server.NewBulletTimeControl(),
				)
				Expect(m.AddMatch(ongoingMatch)).To(Succeed())
			})
			It("returns an error", func() {
				Expect(m.AddMatch(match)).ToNot(Succeed())
			})
		})
		It("adds the match to the active matches", func() {
			Expect(m.AddMatch(match)).To(Succeed())
			Expect(m.GetMatchById(match.Uuid)).To(Equal(match))
			Expect(m.GetMatchByClientKey(match.BlackClientKey)).To(Equal(match))
			Expect(m.GetMatchByClientKey(match.WhiteClientKey)).To(Equal(match))
		})
		It("emits a match created event", func() {
			Expect(m.AddMatch(match)).To(Succeed())

			Eventually(func() int {
				return eventCatcher.EventsByVariantCount(server.MATCH_CREATED)
			}).Should(Equal(1))
			expEvent := server.NewMatchCreatedEvent(match)
			actualEvent := eventCatcher.LastEventByVariant(server.MATCH_CREATED)
			Expect(actualEvent).To(BeEquivalentTo(expEvent))
		})
	})
	Describe("RemoveMatch", func() {
		var match *server.Match
		BeforeEach(func() {
			match = server.NewMatch("client1", "client2", server.NewBulletTimeControl())
		})
		Describe("when the match exists", func() {
			BeforeEach(func() {
				Expect(m.AddMatch(match)).To(Succeed())
			})
			It("removes the match from the active matches", func() {
				Expect(m.RemoveMatch(match)).To(Succeed())
				Expect(m.GetMatchById(match.Uuid)).Error().To(HaveOccurred())
				Expect(m.GetMatchByClientKey(match.WhiteClientKey)).Error().To(HaveOccurred())
				Expect(m.GetMatchByClientKey(match.BlackClientKey)).Error().To(HaveOccurred())
			})
			It("emits a match ended event", func() {
				Expect(m.RemoveMatch(match)).To(Succeed())
				Eventually(func() int {
					return eventCatcher.EventsByVariantCount(server.MATCH_ENDED)
				}).Should(Equal(1))
				expEvent := server.NewMatchEndedEvent(match)
				actualEvent := eventCatcher.LastEventByVariant(server.MATCH_ENDED)
				Expect(actualEvent).To(BeEquivalentTo(expEvent))
			})
		})
		Describe("when the match does not exist", func() {
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
	//	Describe("when the match exists", func() {
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
	//		It("updates the match", func() {
	//			err := m.setMatch(newMatch)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(m.matchByMatchId[newMatch.Uuid]).To(Equal(newMatch))
	//		})
	//		It("emits a match updated event", func() {
	//			err := m.setMatch(newMatch)
	//			Expect(err).ToNot(HaveOccurred())
	//			Expect(el.Events).To(HaveLen(1))
	//			Expect(el.Events[0].Variant()).To(Equal(MATCH_UPDATED))
	//			Expect(el.Events[0].Payload()).To(Equal(newMatch))
	//		})
	//		Describe("when the new match differs by client id", func() {
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
	//		Describe("when the new match differs by time control", func() {
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
	//	Describe("when the match does not exist", func() {
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
	//	var match *Match
	//	var move chess.Move
	//	BeforeEach(func() {
	//		match = NewMatch("client1", "client2", NewBulletTimeControl())
	//		move = chess.Move{
	//			Piece:               chess.WHITE_PAWN,
	//			StartSquare:         &chess.Square{Rank: 2, File: 4},
	//			EndSquare:           &chess.Square{Rank: 4, File: 4},
	//			CapturedPiece:       chess.EMPTY,
	//			KingCheckingSquares: make([]*chess.Square, 0),
	//			PawnUpgradedTo:      chess.EMPTY,
	//		}
	//	})
	//	Describe("when the match doesnt exist", func() {
	//		BeforeEach(func() {
	//			_, getMatchErr := m.GetMatchById(match.Uuid)
	//			Expect(getMatchErr).To(HaveOccurred())
	//		})
	//		It("returns an error", func() {
	//			err := m.ExecuteMove(match.Uuid, &move)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})
	//	Describe("when the match exists", func() {
	//		BeforeEach(func() {
	//			match = NewMatch("client1", "client2", NewBulletTimeControl())
	//			addMatchErr := m.AddMatch(match)
	//			Expect(addMatchErr).ToNot(HaveOccurred())
	//			Expect(m.matchByMatchId).To(HaveKey(match.Uuid))
	//		})
	//		Describe("when the move is illegal", func() {
	//			BeforeEach(func() {
	//				move = chess.Move{chess.WHITE_PAWN, &chess.Square{8, 8}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
	//			})
	//			It("returns an error", func() {
	//				err := m.ExecuteMove(match.Uuid, &move)
	//				Expect(err).To(HaveOccurred())
	//			})
	//		})
	//		It("Updates the match", func() {
	//			err := m.ExecuteMove(match.Uuid, &move)
	//			Expect(err).ToNot(HaveOccurred())
	//			expBoard := chess.GetBoardFromMove(match.Board, &move)
	//			Expect(m.matchByMatchId[match.Uuid].Board).To(Equal(expBoard))
	//		})
	//	})
	//})
})
