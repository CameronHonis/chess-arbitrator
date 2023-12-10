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
				oldMatch.BlackClientKey = "client3"
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
			Expect(mm.matchIdByClientId[match.WhiteClientKey]).To(Equal(match.Uuid))
			Expect(mm.matchIdByClientId[match.BlackClientKey]).To(Equal(match.Uuid))
		})
		It("subscribes the players to the match topic", func() {
			messageTopic := MessageTopic(fmt.Sprintf("match-%s", match.Uuid))
			err := mm.AddMatch(match)
			Expect(err).ToNot(HaveOccurred())
			whiteSubbedTopics := GetSubscriptionManager().GetSubbedTopics(match.WhiteClientKey)
			Expect(whiteSubbedTopics.Flatten()).To(ContainElement(messageTopic))
			blackSubbedTopics := GetSubscriptionManager().GetSubbedTopics(match.BlackClientKey)
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
				mm.matchIdByClientId[match.WhiteClientKey] = match.Uuid
				mm.matchIdByClientId[match.BlackClientKey] = match.Uuid
			})
			It("removes the match from the active matches", func() {
				err := mm.RemoveMatch(match)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.matchByMatchId).ToNot(HaveKey(match.Uuid))
				Expect(mm.matchIdByClientId).ToNot(HaveKey(match.WhiteClientKey))
				Expect(mm.matchIdByClientId).ToNot(HaveKey(match.BlackClientKey))
			})
			It("unsubscribes the players from the match topic", func() {
				messageTopic := MessageTopic(fmt.Sprintf("match-%s", match.Uuid))
				err := mm.RemoveMatch(match)
				Expect(err).ToNot(HaveOccurred())
				whiteSubbedTopics := GetSubscriptionManager().GetSubbedTopics(match.WhiteClientKey)
				Expect(whiteSubbedTopics.Flatten()).ToNot(ContainElement(messageTopic))
				blackSubbedTopics := GetSubscriptionManager().GetSubbedTopics(match.BlackClientKey)
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
			newMatch.WhiteClientKey = "client1"
			newMatch.BlackClientKey = "client2"
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
				prevMatch.WhiteClientKey = "client1"
				prevMatch.BlackClientKey = "client2"
				prevMatch.Uuid = newMatch.Uuid
				mm.matchByMatchId[prevMatch.Uuid] = prevMatch
				mm.matchIdByClientId[prevMatch.WhiteClientKey] = prevMatch.Uuid
				mm.matchIdByClientId[prevMatch.BlackClientKey] = prevMatch.Uuid
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
					newMatch.WhiteClientKey = "other-client1"
					newMatch.BlackClientKey = "client2"
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
					newMatch.WhiteClientKey = "client1"
					newMatch.BlackClientKey = "client2"
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

	FDescribe("StageMatchFromChallenge", func() {
		var challenge Challenge
		BeforeEach(func() {
			challenge = Challenge{
				ChallengerKey:     "client1",
				ChallengedKey:     "client2",
				IsChallengerWhite: true,
				IsChallengerBlack: false,
				TimeControl:       NewBulletTimeControl(),
			}
		})
		When("neither player is in a match", func() {
			When("the challenge specifies that the challenger is white", func() {
				It("stages a match with the challenger as white", func() {
					match, err := mm.StageMatchFromChallenge(&challenge)
					Expect(err).ToNot(HaveOccurred())
					Expect(mm.stagedMatchById[match.Uuid]).To(Equal(match))
					Expect(match.Uuid).ToNot(BeEmpty())
					Expect(match.WhiteClientKey).To(Equal(challenge.ChallengerKey))
					Expect(match.BlackClientKey).To(Equal(challenge.ChallengedKey))
					Expect(match.TimeControl).To(Equal(challenge.TimeControl))
					Expect(match.WhiteTimeRemainingSec).To(Equal(challenge.TimeControl.InitialTimeSec))
					Expect(match.BlackTimeRemainingSec).To(Equal(challenge.TimeControl.InitialTimeSec))
				})
			})
			When("the challenge specifies that the challenger is black", func() {
				BeforeEach(func() {
					challenge.IsChallengerWhite = false
					challenge.IsChallengerBlack = true
				})
				It("stages a match with the challenger as black", func() {
					match, err := mm.StageMatchFromChallenge(&challenge)
					Expect(err).ToNot(HaveOccurred())
					Expect(mm.stagedMatchById[match.Uuid]).To(Equal(match))
					Expect(match.Uuid).ToNot(BeEmpty())
					Expect(match.WhiteClientKey).To(Equal(challenge.ChallengedKey))
					Expect(match.BlackClientKey).To(Equal(challenge.ChallengerKey))
					Expect(match.TimeControl).To(Equal(challenge.TimeControl))
					Expect(match.WhiteTimeRemainingSec).To(Equal(challenge.TimeControl.InitialTimeSec))
					Expect(match.BlackTimeRemainingSec).To(Equal(challenge.TimeControl.InitialTimeSec))
				})
			})
			When("the challenge does not specify player colors", func() {
				BeforeEach(func() {
					challenge.IsChallengerWhite = false
					challenge.IsChallengerBlack = false
				})
				It("randomly assigns player colors", func() {
					memo := make(map[bool]bool)
					for i := 0; i < 100; i++ {
						mm.stagedMatchById = make(map[string]*Match)
						match, err := mm.StageMatchFromChallenge(&challenge)
						Expect(err).ToNot(HaveOccurred())
						memo[match.WhiteClientKey == challenge.ChallengerKey] = true
						if _, ok := memo[match.WhiteClientKey != challenge.ChallengerKey]; ok {
							break
						}
					}
					Expect(len(memo)).To(Equal(2))
				})
			})
			When("the challenger is already in a staged match", func() {
				BeforeEach(func() {
					otherChallenge := Challenge{
						ChallengerKey:     "client1",
						ChallengedKey:     "client3",
					}
					_, err := mm.StageMatchFromChallenge(&otherChallenge)
					Expect(err).ToNot(HaveOccurred())
					Expect(mm.stagedMatchById).To(HaveLen(1)
				})
				It("stages another match", func() {

				})
			})
			When("the challenged is already in a staged match", func() {
				BeforeEach(func() {

				})
				It("stages another match", func() {

				})
			})
		})
		When("the challenger is in a match", func() {
			BeforeEach(func() {

				match, stagingErr := mm.StageMatchFromChallenge(challenge)
			})
			It("returns an error", func() {
				err := mm.StageMatchFromChallenge(match)
				Expect(err).To(HaveOccurred())
			})
		})
		When("the match is already in progress", func() {
			It("returns an error", func() {
				_ = mm.AddMatch(match)
				err := mm.StageMatchFromChallenge(match)
				Expect(err).To(HaveOccurred())
			})
		})
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
