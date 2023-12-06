package server

import (
	"encoding/json"
	"fmt"
	"github.com/CameronHonis/chess"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

type MockUserClientsManager struct {
	UserClientsManagerI
	MessagesBroadcasted []*Message
}

func (mucm *MockUserClientsManager) BroadcastMessage(message *Message) {
	mucm.MessagesBroadcasted = append(mucm.MessagesBroadcasted, message)
}

var _ = Describe("MatchManager", func() {
	var mm *MatchManager
	BeforeEach(func() {
		matchManager = nil
		userClientsManager = nil
		mm = GetMatchManager()
	})
	Describe("AddMatch", func() {
		var match *Match
		BeforeEach(func() {
			match = NewMatch("client1", "client2", NewBulletTimeControl())
		})
		Describe("when one of the players in the proposed match is already in a match", func() {
			BeforeEach(func() {
				mm.matchIdByClientId["client1"] = "match1"
				oldMatch := *match
				oldMatch.BlackClientId = "client3"
				mm.matchByMatchId["match1"] = &oldMatch
			})
			It("returns an error", func() {
				err := mm.AddMatch(match)
				Expect(err).To(HaveOccurred())
			})
		})
		It("adds the match to the active matches", func() {
			err := mm.AddMatch(match)
			Expect(err).ToNot(HaveOccurred())
			Expect(mm.matchByMatchId[match.Uuid]).To(Equal(match))
			Expect(mm.matchIdByClientId[match.WhiteClientId]).To(Equal(match.Uuid))
			Expect(mm.matchIdByClientId[match.BlackClientId]).To(Equal(match.Uuid))
		})
		It("subscribes the players to the match topic", func() {
			messageTopic := MessageTopic(fmt.Sprintf("match-%s", match.Uuid))
			err := mm.AddMatch(match)
			Expect(err).ToNot(HaveOccurred())
			whiteSubbedTopics := GetUserClientsManager().GetSubscribedTopics(match.WhiteClientId)
			Expect(whiteSubbedTopics.Flatten()).To(ContainElement(messageTopic))
			blackSubbedTopics := GetUserClientsManager().GetSubscribedTopics(match.BlackClientId)
			Expect(blackSubbedTopics.Flatten()).To(ContainElement(messageTopic))
		})
	})
	Describe("RemoveMatch", func() {
		var match *Match
		BeforeEach(func() {
			match = NewMatch("client1", "client2", NewBulletTimeControl())
		})
		Describe("when the match exists", func() {
			BeforeEach(func() {
				mm.matchByMatchId[match.Uuid] = match
				mm.matchIdByClientId[match.WhiteClientId] = match.Uuid
				mm.matchIdByClientId[match.BlackClientId] = match.Uuid
			})
			It("removes the match from the active matches", func() {
				err := mm.RemoveMatch(match)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.matchByMatchId).ToNot(HaveKey(match.Uuid))
				Expect(mm.matchIdByClientId).ToNot(HaveKey(match.WhiteClientId))
				Expect(mm.matchIdByClientId).ToNot(HaveKey(match.BlackClientId))
			})
			It("unsubscribes the players from the match topic", func() {
				messageTopic := MessageTopic(fmt.Sprintf("match-%s", match.Uuid))
				err := mm.RemoveMatch(match)
				Expect(err).ToNot(HaveOccurred())
				whiteSubbedTopics := GetUserClientsManager().GetSubscribedTopics(match.WhiteClientId)
				Expect(whiteSubbedTopics.Flatten()).ToNot(ContainElement(messageTopic))
				blackSubbedTopics := GetUserClientsManager().GetSubscribedTopics(match.BlackClientId)
				Expect(blackSubbedTopics.Flatten()).ToNot(ContainElement(messageTopic))
			})
		})
		Describe("when the match does not exist", func() {
			It("returns an error", func() {
				err := mm.RemoveMatch(match)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("SetMatch", func() {
		var newMatch *Match
		BeforeEach(func() {
			newMatch = NewMatch("client1", "client2", NewBulletTimeControl())
			newMatch.WhiteClientId = "client1"
			newMatch.BlackClientId = "client2"
			move := chess.Move{chess.WHITE_PAWN, &chess.Square{2, 4}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
			newBoard := chess.GetBoardFromMove(newMatch.Board, &move)
			newTime := time.Now().Add(time.Second * -10)
			newMatch.Board = newBoard
			newMatch.LastMoveTime = &newTime
		})
		Describe("when the match exists", func() {
			var prevMatch *Match
			BeforeEach(func() {
				prevMatch = NewMatch("client1", "client2", NewBulletTimeControl())
				prevMatch.WhiteClientId = "client1"
				prevMatch.BlackClientId = "client2"
				prevMatch.Uuid = newMatch.Uuid
				mm.matchByMatchId[prevMatch.Uuid] = prevMatch
				mm.matchIdByClientId[prevMatch.WhiteClientId] = prevMatch.Uuid
				mm.matchIdByClientId[prevMatch.BlackClientId] = prevMatch.Uuid
				mm.userClientsManager = &MockUserClientsManager{}
			})
			It("updates the match", func() {
				err := mm.SetMatch(newMatch)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.matchByMatchId[newMatch.Uuid]).To(Equal(newMatch))
			})
			It("broadcasts the match update on the match topic", func() {
				err := mm.SetMatch(newMatch)
				Expect(err).ToNot(HaveOccurred())
				mucm := mm.userClientsManager.(*MockUserClientsManager)
				Expect(mucm.MessagesBroadcasted).To(HaveLen(1))
				messageTopic := MessageTopic(fmt.Sprintf("match-%s", newMatch.Uuid))
				messageJson, _ := mucm.MessagesBroadcasted[0].Marshal()
				expMessageJson, _ := json.Marshal(&Message{
					Topic:       messageTopic,
					ContentType: CONTENT_TYPE_MATCH_UPDATE,
					Content: &MatchUpdateMessageContent{
						Match: newMatch,
					},
				})
				Expect(messageJson).To(Equal(expMessageJson))
			})
			Describe("when the new match differs by client id", func() {
				BeforeEach(func() {
					newMatch = NewMatch("other-client1", "client2", NewBulletTimeControl())
					newMatch.WhiteClientId = "other-client1"
					newMatch.BlackClientId = "client2"
					newMatch.Uuid = prevMatch.Uuid
				})
				It("returns an error", func() {
					err := mm.SetMatch(newMatch)
					Expect(err).To(HaveOccurred())
				})
			})
			Describe("when the new match differs by time control", func() {
				BeforeEach(func() {
					newMatch = NewMatch("client1", "client2", NewRapidTimeControl())
					newMatch.WhiteClientId = "client1"
					newMatch.BlackClientId = "client2"
					newMatch.Uuid = prevMatch.Uuid
				})
				It("returns an error", func() {
					err := mm.SetMatch(newMatch)
					Expect(err).To(HaveOccurred())
				})
			})
		})
		Describe("when the match does not exist", func() {
			BeforeEach(func() {
				Expect(mm.matchByMatchId).ToNot(HaveKey(newMatch.Uuid))
			})
			It("returns an error", func() {
				err := mm.SetMatch(newMatch)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("StageMatch", func() {
	})
	Describe("UnstageMatch", func() {

	})
	Describe("AddMatchFromStaged", func() {

	})

	Describe("ChallengeClient", func() {

	})
	Describe("ExecuteMove", func() {
		var match *Match
		var move chess.Move
		BeforeEach(func() {
			match = NewMatch("client1", "client2", NewBulletTimeControl())
			addMatchErr := mm.AddMatch(match)
			Expect(addMatchErr).ToNot(HaveOccurred())
			Expect(mm.matchByMatchId).To(HaveKey(match.Uuid))
			move = chess.Move{chess.WHITE_PAWN, &chess.Square{2, 4}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
		})
		Describe("when the match doesnt exist", func() {
			BeforeEach(func() {
				mm.matchIdByClientId = make(map[string]string)
				mm.matchByMatchId = make(map[string]*Match)
				_, getMatchErr := mm.GetMatchById(match.Uuid)
				Expect(getMatchErr).To(HaveOccurred())
			})
			It("returns an error", func() {
				err := mm.ExecuteMove(match.Uuid, &move)
				Expect(err).To(HaveOccurred())
			})
		})
		Describe("when the move is illegal", func() {
			BeforeEach(func() {
				move = chess.Move{chess.WHITE_PAWN, &chess.Square{8, 8}, &chess.Square{4, 4}, chess.EMPTY, make([]*chess.Square, 0), chess.EMPTY}
			})
			It("returns an error", func() {
				err := mm.ExecuteMove(match.Uuid, &move)
				Expect(err).To(HaveOccurred())
			})
		})
		It("Updates the match", func() {
			err := mm.ExecuteMove(match.Uuid, &move)
			Expect(err).ToNot(HaveOccurred())
			expBoard := chess.GetBoardFromMove(match.Board, &move)
			Expect(mm.matchByMatchId[match.Uuid].Board).To(Equal(expBoard))
		})
	})
})
