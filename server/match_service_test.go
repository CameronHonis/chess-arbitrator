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

var _ = Describe("MatchService", func() {
	var mm *MatchService
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

	Describe("StageMatchFromChallenge", func() {
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
					Expect(match.WhiteTimeRemainingSec).To(Equal(float64(challenge.TimeControl.InitialTimeSec)))
					Expect(match.BlackTimeRemainingSec).To(Equal(float64(challenge.TimeControl.InitialTimeSec)))
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
					stagedMatchIds := mm.stagedMatchIdsByClientId[challenge.ChallengerKey]
					Expect(stagedMatchIds).ToNot(BeNil())
					Expect(stagedMatchIds.Flatten()).To(HaveLen(1))
					Expect(stagedMatchIds.Flatten()).To(ContainElement(match.Uuid))
					Expect(match.Uuid).ToNot(BeEmpty())
					Expect(match.WhiteClientKey).To(Equal(challenge.ChallengedKey))
					Expect(match.BlackClientKey).To(Equal(challenge.ChallengerKey))
					Expect(match.TimeControl).To(Equal(challenge.TimeControl))
					Expect(match.WhiteTimeRemainingSec).To(Equal(float64(challenge.TimeControl.InitialTimeSec)))
					Expect(match.BlackTimeRemainingSec).To(Equal(float64(challenge.TimeControl.InitialTimeSec)))
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
		})
		When("the challenger is already in a staged match", func() {
			BeforeEach(func() {
				otherChallenge := Challenge{
					ChallengerKey: "client1",
					ChallengedKey: "client3",
					TimeControl:   NewBulletTimeControl(),
				}
				_, err := mm.StageMatchFromChallenge(&otherChallenge)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.stagedMatchById).To(HaveLen(1))
				stagedMatchIds := mm.stagedMatchIdsByClientId[otherChallenge.ChallengerKey]
				Expect(stagedMatchIds.Flatten()).To(HaveLen(1))
			})
			It("stages another match", func() {
				match, err := mm.StageMatchFromChallenge(&challenge)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.stagedMatchById).To(HaveLen(2))
				stagedMatch := mm.stagedMatchById[match.Uuid]
				Expect(stagedMatch).To(Equal(match))
			})
		})
		When("the challenged is already in a staged match", func() {
			BeforeEach(func() {
				otherChallenge := Challenge{
					ChallengerKey: "client2",
					ChallengedKey: "client3",
					TimeControl:   NewBulletTimeControl(),
				}
				_, err := mm.StageMatchFromChallenge(&otherChallenge)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.stagedMatchById).To(HaveLen(1))
				stagedMatchIds := mm.stagedMatchIdsByClientId[otherChallenge.ChallengedKey]
				Expect(stagedMatchIds.Flatten()).To(HaveLen(1))
			})
			It("stages another match", func() {
				match, err := mm.StageMatchFromChallenge(&challenge)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.stagedMatchById).To(HaveLen(2))
				stagedMatch := mm.stagedMatchById[match.Uuid]
				Expect(stagedMatch).To(Equal(match))
			})
		})
		When("the challenger is already in a match", func() {
			BeforeEach(func() {
				addMatchErr := mm.AddMatch(NewMatch("client1", "client3", NewBulletTimeControl()))
				Expect(addMatchErr).ToNot(HaveOccurred())
				Expect(mm.matchByMatchId).To(HaveLen(1))
				challengedClientMatchId := mm.matchIdByClientId["client1"]
				Expect(challengedClientMatchId).ToNot(BeNil())
			})
			It("returns an error", func() {
				_, err := mm.StageMatchFromChallenge(&challenge)
				Expect(err).To(HaveOccurred())
			})
		})
		When("the challenged is already in a match", func() {
			BeforeEach(func() {
				addMatchErr := mm.AddMatch(NewMatch("client2", "client3", NewBulletTimeControl()))
				Expect(addMatchErr).ToNot(HaveOccurred())
				Expect(mm.matchByMatchId).To(HaveLen(1))
				challengedClientMatchId := mm.matchIdByClientId["client2"]
				Expect(challengedClientMatchId).ToNot(BeNil())
			})
			It("returns an error", func() {
				_, err := mm.StageMatchFromChallenge(&challenge)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("UnstageMatch", func() {
		var matchId string
		When("the match is staged", func() {
			BeforeEach(func() {
				stagedMatch, err := mm.StageMatchFromChallenge(&Challenge{
					ChallengerKey: "client1",
					ChallengedKey: "client2",
					TimeControl:   NewBulletTimeControl(),
				})
				matchId = stagedMatch.Uuid
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.stagedMatchById).To(HaveKey(matchId))
				challengerStagedMatchIds, _ := mm.stagedMatchIdsByClientId["client1"]
				Expect(challengerStagedMatchIds.Flatten()).To(ContainElement(matchId))
				challengedStagedMatchIds, _ := mm.stagedMatchIdsByClientId["client2"]
				Expect(challengedStagedMatchIds.Flatten()).To(ContainElement(matchId))
				Expect(mm)
			})
			It("removes the match from the staged matches", func() {
				err := mm.UnstageMatch(matchId)
				Expect(err).ToNot(HaveOccurred())
				Expect(mm.stagedMatchById).ToNot(HaveKey(matchId))
			})
			It("unlinks the client ids from the staged match", func() {
				err := mm.UnstageMatch(matchId)
				Expect(err).ToNot(HaveOccurred())
				challengerStagedMatchIds, _ := mm.stagedMatchIdsByClientId["client1"]
				Expect(challengerStagedMatchIds.Flatten()).ToNot(ContainElement(matchId))
				challengedStagedMatchIds, _ := mm.stagedMatchIdsByClientId["client2"]
				Expect(challengedStagedMatchIds.Flatten()).ToNot(ContainElement(matchId))
			})
		})
		When("no staged match exists with the given id", func() {
			It("returns an error", func() {
				err := mm.UnstageMatch("asdf")
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("AddMatchFromStaged", func() {
		When("the staged match does not exist", func() {
			It("returns an error", func() {
				err := mm.AddMatchFromStaged("asdf")
				Expect(err).To(HaveOccurred())
			})
		})
		When("the staged match exists", func() {
			var stagedMatchId string
			BeforeEach(func() {
				stagedMatch, stageMatchErr := mm.StageMatchFromChallenge(&Challenge{
					ChallengerKey: "client1",
					ChallengedKey: "client2",
					TimeControl:   NewBulletTimeControl(),
				})
				Expect(stageMatchErr).ToNot(HaveOccurred())
				stagedMatchId = stagedMatch.Uuid
				fetchedStagedMatch, getStagedMatchErr := mm.GetStagedMatchById(stagedMatchId)
				Expect(getStagedMatchErr).ToNot(HaveOccurred())
				Expect(fetchedStagedMatch).To(Equal(stagedMatch))
			})
			It("adds the match to the active matches", func() {
				err := mm.AddMatchFromStaged(stagedMatchId)
				Expect(err).ToNot(HaveOccurred())
				match, getMatchErr := mm.GetMatchById(stagedMatchId)
				Expect(getMatchErr).ToNot(HaveOccurred())
				Expect(match).ToNot(BeNil())
			})
			It("removes the staged match", func() {
				err := mm.AddMatchFromStaged(stagedMatchId)
				Expect(err).ToNot(HaveOccurred())
				stagedMatch, getStagedMatchErr := mm.GetStagedMatchById(stagedMatchId)
				Expect(getStagedMatchErr).To(HaveOccurred())
				Expect(stagedMatch).To(BeNil())
			})
			When("the challenger has other staged matches", func() {
				var otherStagedMatchId string
				BeforeEach(func() {
					otherStagedMatch, stageMatchErr := mm.StageMatchFromChallenge(&Challenge{
						ChallengerKey: "client1",
						ChallengedKey: "client3",
						TimeControl:   NewBulletTimeControl(),
					})
					Expect(stageMatchErr).ToNot(HaveOccurred())
					Expect(otherStagedMatch).ToNot(BeNil())
					otherStagedMatchId = otherStagedMatch.Uuid
					challengerStagedMatches, getChallengeStagedMatchesErr := mm.GetStagedMatchesByClientKey("client1")
					Expect(getChallengeStagedMatchesErr).ToNot(HaveOccurred())
					Expect(challengerStagedMatches).To(HaveLen(2))
					fetchedOtherStagedMatch, _ := mm.GetStagedMatchById(otherStagedMatchId)
					Expect(fetchedOtherStagedMatch).To(Equal(otherStagedMatch))
				})
				It("removes the other staged matches", func() {
					err := mm.AddMatchFromStaged(stagedMatchId)
					Expect(err).ToNot(HaveOccurred())
					otherStagedMatch, getOtherStagedMatchErr := mm.GetStagedMatchById(otherStagedMatchId)
					Expect(getOtherStagedMatchErr).To(HaveOccurred())
					Expect(otherStagedMatch).To(BeNil())
					challengerStagedMatches, _ := mm.GetStagedMatchesByClientKey("client1")
					Expect(challengerStagedMatches).To(HaveLen(0))
				})
			})
		})
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
