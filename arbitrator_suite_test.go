package main_test

import (
	"github.com/CameronHonis/chess-arbitrator/app"
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
var clientConn *websocket.Conn
var msgQueue *MsgQueue

var _ = BeforeSuite(func() {
	appService = app.BuildServices()
	appService.Start()
	msgQueue = newMsgQueue()

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
})

var _ = Describe("Auth Workflow", func() {
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
	})
})
