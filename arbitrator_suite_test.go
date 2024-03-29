package main_test

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/app"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"os"
	"sync"
	"testing"
	"time"
)

const PRINT_INBOUND_MSGS = false
const PRINT_OUTBOUND_MSGS = false

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

func connectClient(msgQueue *MsgQueue, clientName string, shouldRequestAuth bool) *websocket.Conn {
	var clientConn *websocket.Conn
	for i := 0; i < 10; i++ {
		clientConn, _, _ = websocket.DefaultDialer.Dial("ws://localhost:8080", nil)
		if clientConn != nil {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}

	if shouldRequestAuth {
		refreshAuthMsg := &models.Message{
			ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
			Content:     &models.RefreshAuthMessageContent{ExistingAuth: nil},
		}
		sendMsg(clientName, clientConn, "", "", refreshAuthMsg)
	}

	go func() {
		for {
			_, msgBytes, readErr := clientConn.ReadMessage()
			if readErr != nil {
				return
			}
			if PRINT_INBOUND_MSGS {
				fmt.Printf("[client %s] << %s\n", clientName, string(msgBytes))
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

func sendMsg(clientName string, conn *websocket.Conn, pubKey, privKey models.Key, msg *models.Message) {
	msg.SenderKey = pubKey
	msg.PrivateKey = privKey
	msgBytes, marshalErr := msg.Marshal()
	if marshalErr != nil {
		panic(marshalErr)
	}
	if PRINT_OUTBOUND_MSGS {
		fmt.Printf("[client %s] >> %s\n", clientName, string(msgBytes))
	}
	if writeErr := conn.WriteMessage(websocket.TextMessage, msgBytes); writeErr != nil {
		panic(writeErr)
	}
}

func listenForMsgType(msgQueue *MsgQueue, contentType models.ContentType) *models.Message {
	for i := 0; i < 100; i++ {
		msgs := msgQueue.toSlice()
		for _, msg := range msgs {
			if msg.ContentType == contentType {
				return msg
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	panic(fmt.Sprintf("timed out waiting for msg of type %s", contentType))
}

var botClientSecret string
var prevBotClientSecret string
var appService app.AppServiceI
var _ = BeforeSuite(func() {
	appService = app.BuildServices(app.GetMutedLoggerConfig())
	appService.Start()

	botClientSecret = "bot_client_secret"
	prevBotClientSecret = os.Getenv(string(models.SECRET_BOT_CLIENT_SECRET))
	Expect(os.Setenv(string(models.SECRET_BOT_CLIENT_SECRET), botClientSecret)).To(Succeed())
})

var _ = AfterSuite(func() {
	appService.Stop()

	_ = os.Setenv(string(models.SECRET_BOT_CLIENT_SECRET), prevBotClientSecret)
})

var _ = Describe("integration tests", func() {
	Describe("auth", func() {
		var prevAuthKeyMinsToStale string
		var conn *websocket.Conn
		var msgQueue *MsgQueue
		BeforeEach(func() {
			msgQueue = newMsgQueue()
			conn = connectClient(msgQueue, "A", false)
			prevAuthKeyMinsToStale = os.Getenv(string(models.SECRET_AUTH_KEY_MINS_TO_STALE))
			_ = os.Setenv(string(models.SECRET_AUTH_KEY_MINS_TO_STALE), "1")
		})
		AfterEach(func() {
			_ = os.Setenv(string(models.SECRET_AUTH_KEY_MINS_TO_STALE), prevAuthKeyMinsToStale)
		})
		Describe("request refresh auth", func() {
			When("no prior creds exist", func() {
				BeforeEach(func() {
					refreshAuthMsg := &models.Message{
						ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
						Content: &models.RefreshAuthMessageContent{
							ExistingAuth: nil,
						},
					}
					sendMsg("A", conn, "", "", refreshAuthMsg)
				})
				It("replies with fresh creds", func() {
					listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
				})
			})
			When("prior creds do exist", func() {
				var pubKey models.Key
				var priKey models.Key
				BeforeEach(func() {
					refreshAuthMsg := &models.Message{
						ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
						Content: &models.RefreshAuthMessageContent{
							ExistingAuth: nil,
						},
					}
					sendMsg("A", conn, "", "", refreshAuthMsg)
					msg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
					authMsgContent := msg.Content.(*models.AuthMessageContent)
					pubKey = authMsgContent.PublicKey
					priKey = authMsgContent.PrivateKey
					msgQueue.flush()
					Expect(conn.Close()).ToNot(HaveOccurred())

					conn = connectClient(msgQueue, "A", false)
				})
				When("the auth is valid", func() {
					When("the auth is fresh", func() {
						BeforeEach(func() {
							msgQueue.flush()
							refreshAuthMsg := &models.Message{
								ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
								Content: &models.RefreshAuthMessageContent{
									ExistingAuth: &models.AuthMessageContent{
										PublicKey:  pubKey,
										PrivateKey: priKey,
									},
								},
							}
							sendMsg("A", conn, pubKey, priKey, refreshAuthMsg)
						})
						It("mirrors the creds back", func() {
							authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
							authMsgContent := authMsg.Content.(*models.AuthMessageContent)
							Expect(authMsgContent.PublicKey).To(Equal(pubKey))
							Expect(authMsgContent.PrivateKey).To(Equal(priKey))
						})
					})
					When("the auth is stale", func() {
						BeforeEach(func() {
							msgQueue.flush()

							prevAuthKeyMinsToStale = os.Getenv(string(models.SECRET_AUTH_KEY_MINS_TO_STALE))
							_ = os.Setenv(string(models.SECRET_AUTH_KEY_MINS_TO_STALE), "0.0001")

							time.Sleep(time.Millisecond * 100)

							refreshAuthMsg := &models.Message{
								ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
								Content: &models.RefreshAuthMessageContent{
									ExistingAuth: &models.AuthMessageContent{
										PublicKey:  pubKey,
										PrivateKey: priKey,
									},
								},
							}
							sendMsg("A", conn, pubKey, priKey, refreshAuthMsg)
						})
						AfterEach(func() {
							_ = os.Setenv(string(models.SECRET_AUTH_KEY_MINS_TO_STALE), prevAuthKeyMinsToStale)
						})
						It("replies with the updated private key", func() {
							authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
							authMsgContent := authMsg.Content.(*models.AuthMessageContent)
							Expect(authMsgContent.PublicKey).To(Equal(pubKey))
							Expect(authMsgContent.PrivateKey).ToNot(Equal(priKey))
						})
					})
				})
				When("the auth is invalid", func() {
					BeforeEach(func() {
						msgQueue.flush()
						refreshAuthMsg := &models.Message{
							ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
							Content: &models.RefreshAuthMessageContent{
								ExistingAuth: &models.AuthMessageContent{
									PublicKey:  pubKey,
									PrivateKey: "invalid",
								},
							},
						}
						sendMsg("A", conn, pubKey, priKey, refreshAuthMsg)
					})
					It("replies with a new client creds", func() {
						authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
						authMsgContent := authMsg.Content.(*models.AuthMessageContent)
						Expect(authMsgContent.PublicKey).ToNot(Equal(pubKey))
						Expect(authMsgContent.PrivateKey).ToNot(Equal(priKey))
					})
				})
			})
			When("the client is currently in a match", func() {
				var pubKey models.Key
				var priKey models.Key
				BeforeEach(func() {
					refreshAuthMsg := &models.Message{
						ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
						Content: &models.RefreshAuthMessageContent{
							ExistingAuth: nil,
						},
					}
					sendMsg("A", conn, "", "", refreshAuthMsg)
					authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
					authMsgContent := authMsg.Content.(*models.AuthMessageContent)
					pubKey = authMsgContent.PublicKey
					priKey = authMsgContent.PrivateKey

					msgQueueB := newMsgQueue()
					connB := connectClient(msgQueueB, "B", true)
					authBMsg := listenForMsgType(msgQueueB, models.CONTENT_TYPE_AUTH)
					authMsgContent = authBMsg.Content.(*models.AuthMessageContent)
					pubKeyB := authMsgContent.PublicKey
					priKeyB := authMsgContent.PrivateKey
					msgQueueB.flush()

					sendMsg("A", conn, pubKey, priKey, &models.Message{
						ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST,
						Content: &models.ChallengeRequestMessageContent{
							Challenge: builders.NewChallenge(pubKey, pubKeyB, true, false, builders.NewBlitzTimeControl(), "", true),
						},
					})

					_ = listenForMsgType(msgQueueB, models.CONTENT_TYPE_CHALLENGE_UPDATED)

					msgQueue.flush()
					msgQueueB.flush()
					sendMsg("B", connB, pubKeyB, priKeyB, &models.Message{
						ContentType: models.CONTENT_TYPE_ACCEPT_CHALLENGE,
						Content: &models.AcceptChallengeMessageContent{
							ChallengerClientKey: pubKey,
						},
					})

					_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_MATCH_UPDATED)
					_ = listenForMsgType(msgQueueB, models.CONTENT_TYPE_MATCH_UPDATED)
					msgQueue.flush()
					msgQueueB.flush()

					Expect(conn.Close()).ToNot(HaveOccurred())

					conn = connectClient(msgQueue, "A", false)
				})
				It("responds with the match update", func() {
					refreshAuthMsg := &models.Message{
						ContentType: models.CONTENT_TYPE_REFRESH_AUTH,
						Content: &models.RefreshAuthMessageContent{
							ExistingAuth: &models.AuthMessageContent{
								PublicKey:  pubKey,
								PrivateKey: priKey,
							},
						},
					}
					sendMsg("A", conn, pubKey, priKey, refreshAuthMsg)

					listenForMsgType(msgQueue, models.CONTENT_TYPE_MATCH_UPDATED)
				})
			})
		})
	})

	Describe("auth upgrade", func() {
		var conn *websocket.Conn
		var msgQueue *MsgQueue
		BeforeEach(func() {
			msgQueue = newMsgQueue()
			conn = connectClient(msgQueue, "A", true)
		})
		Describe("request auth upgrade", func() {
			When("the request includes valid auth keys", func() {
				var pubKey models.Key
				var privKey models.Key
				BeforeEach(func() {
					authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
					msgQueue.flush()

					pubKey = authMsg.Content.(*models.AuthMessageContent).PublicKey
					privKey = authMsg.Content.(*models.AuthMessageContent).PrivateKey

					sendMsg("A", conn, pubKey, privKey, &models.Message{
						ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
						Content: &models.UpgradeAuthRequestMessageContent{
							Role:   models.BOT,
							Secret: botClientSecret,
						},
					})

				})
				It("responds with an Upgrade Auth Msg", func() {
					authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED)
					Expect(authMsg).To(PointTo(HaveField(
						"Content", PointTo(MatchAllFields(Fields{
							"UpgradedToRole": BeEquivalentTo(models.BOT),
						})),
					)))
				})
				AfterEach(func() {
					_ = conn.Close()
				})
			})
			When("the request doesn't include auth keys", func() {
				It("responds with an invalid auth msg", func() {
					sendMsg("A", conn, "", "", &models.Message{
						ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
						Content: &models.UpgradeAuthRequestMessageContent{
							Role:   models.BOT,
							Secret: "secret",
						},
					})

					listenForMsgType(msgQueue, models.CONTENT_TYPE_INVALID_AUTH)
				})
			})
		})

	})

	Describe("challenges", func() {
		var conn *websocket.Conn
		var msgQueue *MsgQueue
		BeforeEach(func() {
			msgQueue = newMsgQueue()
			conn = connectClient(msgQueue, "A", true)
		})
		When("client A sends client B a challenge request", func() {
			var clientAPubKey models.Key
			var clientAPrivKey models.Key
			var clientBConn *websocket.Conn
			var clientBMsgQueue *MsgQueue
			var clientBPubKey models.Key
			var clientBPrivKey models.Key
			var challengeAtoB *models.Challenge
			BeforeEach(func() {
				clientBMsgQueue = newMsgQueue()
				clientBConn = connectClient(clientBMsgQueue, "B", true)

				authMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
				msgQueue.flush()
				challengedAuthMsg := listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_AUTH)
				clientBMsgQueue.flush()

				clientAPubKey = authMsg.Content.(*models.AuthMessageContent).PublicKey
				clientAPrivKey = authMsg.Content.(*models.AuthMessageContent).PrivateKey
				clientBPubKey = challengedAuthMsg.Content.(*models.AuthMessageContent).PublicKey
				clientBPrivKey = challengedAuthMsg.Content.(*models.AuthMessageContent).PrivateKey
				challengeAtoB = builders.NewChallenge(
					clientAPubKey,
					clientBPubKey,
					true,
					false,
					builders.NewBlitzTimeControl(),
					"",
					true)
				sendMsg("A", conn, clientAPubKey, clientAPrivKey, &models.Message{
					ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST,
					Content: &models.ChallengeRequestMessageContent{
						Challenge: challengeAtoB,
					},
				})
			})
			It("sends clients A & B the new challenge", func() {
				challengeUpdatedMsgToA := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
				Expect(challengeUpdatedMsgToA).To(PointTo(HaveField(
					"Content", PointTo(HaveField(
						"Challenge", PointTo(MatchAllFields(Fields{
							"Uuid":              Not(BeNil()),
							"ChallengerKey":     Equal(challengeAtoB.ChallengerKey),
							"ChallengedKey":     Equal(challengeAtoB.ChallengedKey),
							"IsChallengerWhite": Equal(challengeAtoB.IsChallengerWhite),
							"IsChallengerBlack": Equal(challengeAtoB.IsChallengerBlack),
							"TimeControl":       Equal(challengeAtoB.TimeControl),
							"BotName":           BeEmpty(),
							"TimeCreated":       PointTo(BeTemporally(">=", time.Now().Add(-1*time.Second))),
							"IsActive":          BeTrue(),
						}))),
					),
				)))

				challengeUpdateMsgToB := listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
				Expect(challengeUpdateMsgToB).To(Equal(challengeUpdatedMsgToA))
			})
			Describe("and the challenger client revokes the challenge", func() {
				BeforeEach(func() {
					_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					msgQueue.flush()
					_ = listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					clientBMsgQueue.flush()

					sendMsg("A", conn, clientAPubKey, clientAPrivKey, &models.Message{
						ContentType: models.CONTENT_TYPE_REVOKE_CHALLENGE,
						Content: &models.RevokeChallengeMessageContent{
							ChallengedClientKey: challengeAtoB.ChallengedKey,
						},
					})
				})
				It("responds to clients A & B with a challenge update msg", func() {
					challengeUpdatedMsgToA := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					Expect(challengeUpdatedMsgToA).To(PointTo(HaveField(
						"Content", PointTo(HaveField(
							"Challenge", PointTo(MatchAllFields(Fields{
								"Uuid":              Not(BeNil()),
								"ChallengerKey":     Equal(challengeAtoB.ChallengerKey),
								"ChallengedKey":     Equal(challengeAtoB.ChallengedKey),
								"IsChallengerWhite": Equal(challengeAtoB.IsChallengerWhite),
								"IsChallengerBlack": Equal(challengeAtoB.IsChallengerBlack),
								"TimeControl":       Equal(challengeAtoB.TimeControl),
								"BotName":           BeEmpty(),
								"TimeCreated":       PointTo(BeTemporally(">=", time.Now().Add(-2*time.Second))),
								"IsActive":          BeFalse(),
							})),
						)),
					)))

					challengeUpdateMsgToB := listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					Expect(challengeUpdateMsgToB).To(Equal(challengeUpdatedMsgToA))
				})
			})
			Describe("and client A disconnects", func() {
				Describe("and client B accepts", func() {
					It("?", func() {

					})
				})
			})
			Describe("and client B accepts", func() {
				BeforeEach(func() {
					_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					msgQueue.flush()
					_ = listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					clientBMsgQueue.flush()

					sendMsg("B", clientBConn, clientBPubKey, clientBPrivKey, &models.Message{
						ContentType: models.CONTENT_TYPE_ACCEPT_CHALLENGE,
						Content: &models.AcceptChallengeMessageContent{
							ChallengerClientKey: challengeAtoB.ChallengerKey,
						},
					})
				})
				It("responds to clients A & B with an inactive challenge msg", func() {
					challengeAcceptedMsgToA := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					Expect(challengeAcceptedMsgToA).To(PointTo(HaveField(
						"Content", PointTo(HaveField(
							"Challenge", PointTo(MatchAllFields(Fields{
								"Uuid":              Not(BeNil()),
								"ChallengerKey":     Equal(challengeAtoB.ChallengerKey),
								"ChallengedKey":     Equal(challengeAtoB.ChallengedKey),
								"IsChallengerWhite": Equal(challengeAtoB.IsChallengerWhite),
								"IsChallengerBlack": Equal(challengeAtoB.IsChallengerBlack),
								"TimeControl":       Equal(challengeAtoB.TimeControl),
								"BotName":           BeEmpty(),
								"TimeCreated":       PointTo(BeTemporally(">=", time.Now().Add(-2*time.Second))),
								"IsActive":          BeFalse(),
							})),
						)),
					)))

					challengeAcceptedMsgToB := listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					Expect(challengeAcceptedMsgToB).To(Equal(challengeAcceptedMsgToA))
				})
				It("responds to clients A & B with a match created msg", func() {
					matchCreatedMsgToA := listenForMsgType(msgQueue, models.CONTENT_TYPE_MATCH_UPDATED)
					Expect(matchCreatedMsgToA).To(PointTo(HaveField(
						"Content", PointTo(HaveField(
							"Match", PointTo(MatchAllFields(Fields{
								"Uuid":                  Not(BeNil()),
								"Board":                 Equal(chess.GetInitBoard()),
								"WhiteClientKey":        Equal(challengeAtoB.ChallengerKey),
								"WhiteTimeRemainingSec": Equal(300.0),
								"BlackClientKey":        Equal(challengeAtoB.ChallengedKey),
								"BlackTimeRemainingSec": Equal(300.0),
								"TimeControl":           Equal(challengeAtoB.TimeControl),
								"BotName":               Equal(""),
								"LastMoveTime":          Ignore(),
								"LastMove":              BeNil(),
								"Result":                Equal(models.MATCH_RESULT_IN_PROGRESS),
							})),
						)),
					)))

					matchCreatedMsgToB := listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_MATCH_UPDATED)
					Expect(matchCreatedMsgToB).To(Equal(matchCreatedMsgToA))
				})
			})
			Describe("and the challenged declines", func() {
				BeforeEach(func() {
					_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					msgQueue.flush()
					_ = listenForMsgType(clientBMsgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					clientBMsgQueue.flush()

					sendMsg("B", clientBConn, clientBPubKey, clientBPrivKey, &models.Message{
						ContentType: models.CONTENT_TYPE_DECLINE_CHALLENGE,
						Content: &models.DeclineChallengeMessageContent{
							ChallengerClientKey: challengeAtoB.ChallengerKey,
						},
					})
				})
				It("responds to both clients with an inactive challenge update", func() {
					challengeUpdatedMsg := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					Expect(challengeUpdatedMsg).To(PointTo(HaveField(
						"Content", PointTo(HaveField(
							"Challenge", PointTo(HaveField(
								"IsActive", BeFalse(),
							)),
						)),
					)))
				})
			})
		})

		When("clients A & B are in a match", func() {
			var pubKeyA models.Key
			var privKeyA models.Key
			var msgQueueB *MsgQueue
			var connB *websocket.Conn
			var pubKeyB models.Key
			var privKeyB models.Key
			BeforeEach(func() {
				msgQueueB = newMsgQueue()
				connB = connectClient(msgQueueB, "B", true)

				authMsgToA := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH)
				pubKeyA = authMsgToA.Content.(*models.AuthMessageContent).PublicKey
				privKeyA = authMsgToA.Content.(*models.AuthMessageContent).PrivateKey

				authMsgToB := listenForMsgType(msgQueueB, models.CONTENT_TYPE_AUTH)
				pubKeyB = authMsgToB.Content.(*models.AuthMessageContent).PublicKey
				privKeyB = authMsgToB.Content.(*models.AuthMessageContent).PrivateKey
				_ = authMsgToB.Content.(*models.AuthMessageContent).PrivateKey

				sendMsg("A", conn, pubKeyA, privKeyA, &models.Message{
					ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST,
					Content: &models.ChallengeRequestMessageContent{
						Challenge: builders.NewChallenge(
							pubKeyA,
							pubKeyB,
							true,
							false,
							builders.NewBlitzTimeControl(),
							"",
							true),
					},
				})

				_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
				msgQueue.flush()
				_ = listenForMsgType(msgQueueB, models.CONTENT_TYPE_CHALLENGE_UPDATED)
				msgQueueB.flush()

				sendMsg("B", connB, pubKeyB, privKeyB, &models.Message{
					ContentType: models.CONTENT_TYPE_ACCEPT_CHALLENGE,
					Content: &models.AcceptChallengeMessageContent{
						ChallengerClientKey: pubKeyA,
					},
				})

				_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
				msgQueue.flush()
				_ = listenForMsgType(msgQueueB, models.CONTENT_TYPE_CHALLENGE_UPDATED)
				msgQueueB.flush()
			})

			Describe("and a third client challenges client A", func() {
				var msgQueueC *MsgQueue
				var pubKeyC models.Key
				BeforeEach(func() {
					msgQueueC = newMsgQueue()
					thirdClientConn := connectClient(msgQueueC, "C", true)
					thirdAuthMsg := listenForMsgType(msgQueueC, models.CONTENT_TYPE_AUTH)
					pubKeyC = thirdAuthMsg.Content.(*models.AuthMessageContent).PublicKey
					thirdClientPrivKey := thirdAuthMsg.Content.(*models.AuthMessageContent).PrivateKey

					sendMsg("C", thirdClientConn, pubKeyC, thirdClientPrivKey, &models.Message{
						ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST,
						Content: &models.ChallengeRequestMessageContent{
							Challenge: builders.NewChallenge(
								pubKeyC,
								pubKeyA,
								true,
								false,
								builders.NewBlitzTimeControl(),
								"",
								true),
						},
					})
				})
				It("responds to client C with a challenge update msg", func() {
					challengeUpdatedMsgToClientC := listenForMsgType(msgQueueC, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					Expect(challengeUpdatedMsgToClientC).To(PointTo(HaveField(
						"Content", PointTo(HaveField(
							"Challenge", PointTo(MatchAllFields(Fields{
								"Uuid":              Not(BeNil()),
								"ChallengerKey":     Equal(pubKeyC),
								"ChallengedKey":     Equal(pubKeyA),
								"IsChallengerWhite": Equal(true),
								"IsChallengerBlack": Equal(false),
								"TimeControl":       Equal(builders.NewBlitzTimeControl()),
								"BotName":           BeEmpty(),
								"TimeCreated":       PointTo(BeTemporally("~", time.Now(), time.Second)),
								"IsActive":          BeTrue(),
							}))),
						))))

					challengeUpdatedMsgToClientA := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
					Expect(challengeUpdatedMsgToClientA).To(Equal(challengeUpdatedMsgToClientC))
				})
				Describe("and client A accepts", func() {
					BeforeEach(func() {
						_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
						msgQueue.flush()
						_ = listenForMsgType(msgQueueC, models.CONTENT_TYPE_CHALLENGE_UPDATED)
						msgQueueC.flush()

						sendMsg("A", conn, pubKeyA, privKeyA, &models.Message{
							ContentType: models.CONTENT_TYPE_ACCEPT_CHALLENGE,
							Content: &models.AcceptChallengeMessageContent{
								ChallengerClientKey: pubKeyC,
							},
						})
					})
					It("responds to clients A & C with an inactive challenge update msg", func() {
						challengeUpdatedMsgToA := listenForMsgType(msgQueue, models.CONTENT_TYPE_CHALLENGE_UPDATED)
						Expect(challengeUpdatedMsgToA).To(PointTo(HaveField(
							"Content", PointTo(HaveField(
								"Challenge", PointTo(MatchFields(IgnoreExtras, Fields{
									"IsActive": BeFalse()})),
							)),
						)))
					})
					It("responds to clients A & C with a match creation failed msg", func() {
						matchCreationFailedMsgToA := listenForMsgType(msgQueue, models.CONTENT_TYPE_MATCH_CREATION_FAILED)
						Expect(matchCreationFailedMsgToA).To(PointTo(HaveField(
							"Content", PointTo(HaveField(
								"Reason", Not(BeEmpty()),
							)),
						)))

						matchCreationFailedMsgToC := listenForMsgType(msgQueueC, models.CONTENT_TYPE_MATCH_CREATION_FAILED)
						Expect(matchCreationFailedMsgToC).To(Equal(matchCreationFailedMsgToA))
					})
				})
			})
			Describe("and client B challenges client A", func() {
				BeforeEach(func() {
					sendMsg("B", connB, pubKeyB, privKeyB, &models.Message{
						ContentType: models.CONTENT_TYPE_CHALLENGE_REQUEST,
						Content: &models.ChallengeRequestMessageContent{
							Challenge: builders.NewChallenge(
								pubKeyB,
								pubKeyA,
								true,
								false,
								builders.NewBlitzTimeControl(),
								"",
								true),
						},
					})
				})
				It("responds to client B with a challenge request failed msg", func() {
					challengeReqFailedMsg := listenForMsgType(msgQueueB, models.CONTENT_TYPE_CHALLENGE_REQUEST_FAILED)
					Expect(challengeReqFailedMsg).To(PointTo(HaveField(
						"Content", PointTo(HaveField(
							"Reason", Not(BeEmpty()),
						)),
					)))
				})
			})
		})
	})

	// test "journeys" below
	Describe("journeys", func() {
		It("allows a bot to join and rejoin", func() {
			msgQueue := newMsgQueue()
			conn := connectClient(msgQueue, "A", true)
			authMsgContent := listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH).Content.(*models.AuthMessageContent)
			pubKey := authMsgContent.PublicKey
			privKey := authMsgContent.PrivateKey

			msgQueue.flush()
			sendMsg("A", conn, pubKey, privKey, &models.Message{
				ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
				Content: &models.UpgradeAuthRequestMessageContent{
					Role:   models.BOT,
					Secret: botClientSecret,
				},
			})

			_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED)
			msgQueue.flush()

			Expect(conn.Close()).To(Succeed())

			conn = connectClient(msgQueue, "A", true)
			authMsgContent = listenForMsgType(msgQueue, models.CONTENT_TYPE_AUTH).Content.(*models.AuthMessageContent)
			msgQueue.flush()
			pubKey = authMsgContent.PublicKey
			privKey = authMsgContent.PrivateKey

			sendMsg("A", conn, pubKey, privKey, &models.Message{
				ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
				Content: &models.UpgradeAuthRequestMessageContent{
					Role:   models.BOT,
					Secret: botClientSecret,
				},
			})
			_ = listenForMsgType(msgQueue, models.CONTENT_TYPE_UPGRADE_AUTH_GRANTED)
		})
	})
})
