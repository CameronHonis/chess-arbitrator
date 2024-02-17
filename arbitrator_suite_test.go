package main_test

import (
	"github.com/CameronHonis/chess-arbitrator/app"
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"sync"
	"testing"
	"time"
)

func TestArbitrator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arbitrator Suite")
}

type MsgQueue struct {
	dump []*models.Message
	ptr  int
	mu   sync.Mutex
}

func newMsgQueue() *MsgQueue {
	return &MsgQueue{
		dump: make([]*models.Message, 0),
	}
}

func (m *MsgQueue) push(msg *models.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.dump = append(m.dump, msg)
}

func (m *MsgQueue) pop() (msg *models.Message, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ptr >= len(m.dump) {
		return nil, false
	}
	msg = m.dump[m.ptr]
	m.ptr++
	return msg, true
}

func (m *MsgQueue) at(idx int) (msg *models.Message, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ptr+idx >= len(m.dump) {
		return nil, false
	}
	return m.dump[m.ptr+idx], true
}

func (m *MsgQueue) flush() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ptr = 0
	m.dump = make([]*models.Message, 0)
}

func (m *MsgQueue) toSlice() []*models.Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.dump[m.ptr:]
}

var appService app.AppServiceI

var _ = BeforeSuite(func() {
	appService = app.BuildServices()
	appService.Start()
})

func connectClient(msgQueue *MsgQueue) *websocket.Conn {
	var clientConn *websocket.Conn
	for i := 0; i < 10; i++ {
		clientConn, _, _ = websocket.DefaultDialer.Dial("ws://localhost:8080", nil)
		if clientConn != nil {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}

	go func() {
		for {
			_, msgBytes, readErr := clientConn.ReadMessage()
			if readErr != nil {
				panic(readErr)
			}

			msg, unmarshalErr := models.UnmarshalToMessage(msgBytes)
			if unmarshalErr != nil {
				panic(unmarshalErr)
			}

			msgQueue.push(msg)
		}
	}()

	return clientConn
}

func sendMsg(conn *websocket.Conn, pubKey, privKey models.Key, msg *models.Message) {
	msg.SenderKey = pubKey
	msg.PrivateKey = privKey
	msgBytes, marshalErr := msg.Marshal()
	if marshalErr != nil {
		panic(marshalErr)
	}
	if writeErr := conn.WriteMessage(websocket.TextMessage, msgBytes); writeErr != nil {
		panic(writeErr)
	}
}

func listenForMsgType(msgQueue *MsgQueue, contentType models.ContentType) *models.Message {
	var matchedMsg *models.Message
	Eventually(func() *models.Message {
		msgs := msgQueue.toSlice()
		for _, msg := range msgs {
			if msg.ContentType == contentType {
				matchedMsg = msg
				return msg
			}
		}
		return nil
	}).ShouldNot(BeNil())

	return matchedMsg
}

var _ = Describe("Workflows", func() {
	var conn *websocket.Conn
	var msgQueue *MsgQueue
	BeforeEach(func() {
		msgQueue = newMsgQueue()
		conn = connectClient(msgQueue)
	})
	Describe("receives basic auth upon connection", func() {
		It("responds with an Auth Msg", func() {
			Eventually(func() []*models.Message {
				msgs := msgQueue.toSlice()
				return msgs
			}).Should(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"ContentType": Equal(models.CONTENT_TYPE_AUTH),
				"Content": PointTo(MatchFields(IgnoreExtras, Fields{
					"PublicKey":  Not(BeZero()),
					"PrivateKey": Not(BeZero()),
				})),
			}))))
			msgQueue.flush()
		})
	})

	Describe("request auth upgrade", func() {
		It("responds with an Upgrade Auth Msg", func() {
			authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
			msgQueue.flush()

			pubKey := authMsg.Content.(*models.AuthMessageContent).PublicKey
			privKey := authMsg.Content.(*models.AuthMessageContent).PrivateKey
			secret, _ := (&auth.AuthenticationService{}).GetSecret(models.BOT)
			sendMsg(conn, pubKey, privKey, &models.Message{
				ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
				Content: &models.UpgradeAuthRequestMessageContent{
					Role:   models.BOT,
					Secret: secret,
				},
			})

			Eventually(func() []*models.Message {
				msgs := msgQueue.toSlice()
				return msgs
			}).Should(ContainElement(PointTo(MatchFields(IgnoreExtras, Fields{
				"ContentType": Equal(models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED),
				"Content": PointTo(MatchAllFields(Fields{
					"UpgradedToRole": BeEquivalentTo(models.BOT),
				})),
			}))))
		})
	})

	FDescribe("send challenge", func() {
		//var challengedConn *websocket.Conn
		var challengedMsgQueue *MsgQueue
		var challengerPubKey models.Key
		var challengerPrivKey models.Key
		var challengedPubKey models.Key
		var challenge *models.Challenge
		BeforeEach(func() {
			challengedMsgQueue = newMsgQueue()
			_ = connectClient(challengedMsgQueue)

			authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
			msgQueue.flush()
			challengedAuthMsg := listenForMsgType(challengedMsgQueue, models.CONTENT_TYPE_AUTH)
			challengedMsgQueue.flush()

			challengerPubKey = authMsg.Content.(*models.AuthMessageContent).PublicKey
			challengerPrivKey = authMsg.Content.(*models.AuthMessageContent).PrivateKey
			challengedPubKey = challengedAuthMsg.Content.(*models.AuthMessageContent).PublicKey
			challenge = builders.NewChallenge(
				challengerPubKey,
				challengedPubKey,
				true,
				false,
				builders.NewBlitzTimeControl(),
				"",
				true)
			sendMsg(conn, challengerPubKey, challengerPrivKey, &models.Message{
				ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST,
				Content: &models.ChallengeRequestMessageContent{
					Challenge: challenge,
				},
			})
		})
		It("sends both the challenger and challenged client the new challenge", func() {
			challengeUpdatedMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
			Expect(challengeUpdatedMsg).To(PointTo(HaveField(
				"Content", PointTo(HaveField(
					"Challenge", PointTo(MatchAllFields(Fields{
						"Uuid":              Ignore(),
						"ChallengerKey":     Equal(challenge.ChallengerKey),
						"ChallengedKey":     Equal(challenge.ChallengedKey),
						"IsChallengerWhite": Equal(challenge.IsChallengerWhite),
						"IsChallengerBlack": Equal(challenge.IsChallengerBlack),
						"TimeControl":       Equal(challenge.TimeControl),
						"BotName":           BeEmpty(),
						"TimeCreated":       PointTo(BeTemporally(">=", time.Now().Add(-1*time.Second))),
						"IsActive":          BeTrue(),
					}))),
				),
			)))

			challengedChallengeUpdateMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
			Expect(challengedChallengeUpdateMsg).To(Equal(challengeUpdatedMsg))
		})
		Describe("and the challenger client revokes the challenge", func() {
			It("responds to both clients with a challenge update msg", func() {

			})
		})
		Describe("and the challenger client disconnects", func() {
			Describe("and the challenged client accepts", func() {
				It("?", func() {

				})
			})
		})
		Describe("and the challenged client accepts", func() {
			It("responds to both clients with a challenge accepted msg", func() {

			})
			It("responds to both clients with a match created msg", func() {

			})
		})
		Describe("and the challenged declines", func() {
			It("responds to both clients with a challenge removed msg", func() {

			})
		})
	})
})
