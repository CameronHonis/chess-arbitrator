package matchmaking

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
	"math"
	"sync"
	"time"
)

type MMPoolNode struct {
	next          *MMPoolNode
	prev          *MMPoolNode
	clientProfile *models.ClientProfile
	timeControl   *models.TimeControl
	timeJoined    int64
}

func NewMMPoolNode(profile *models.ClientProfile, timeControl *models.TimeControl) *MMPoolNode {
	return &MMPoolNode{
		clientProfile: profile,
		timeControl:   timeControl,
		timeJoined:    time.Now().Unix(),
	}
}

func (n *MMPoolNode) Next() *MMPoolNode {
	return n.next
}

func (n *MMPoolNode) Prev() *MMPoolNode {
	return n.prev
}

func (n *MMPoolNode) ClientProfile() *models.ClientProfile {
	return n.clientProfile
}

type MatchmakingPool struct {
	// doubly linked list of client profiles, sorted by queue time (oldest to newest)
	// prioritizes add/remove speed
	head *MMPoolNode
	tail *MMPoolNode
	// map to allow for O(1) lookup time of nodes by client key
	nodeByClientKey map[models.Key]*MMPoolNode
	mu              sync.Mutex
}

func NewMatchmakingPool() *MatchmakingPool {
	return &MatchmakingPool{
		nodeByClientKey: make(map[models.Key]*MMPoolNode),
		mu:              sync.Mutex{},
	}
}

func (mmp *MatchmakingPool) Head() *MMPoolNode {
	mmp.mu.Lock()
	defer mmp.mu.Unlock()
	return mmp.head
}

func (mmp *MatchmakingPool) Tail() *MMPoolNode {
	mmp.mu.Lock()
	defer mmp.mu.Unlock()
	return mmp.tail
}

func (mmp *MatchmakingPool) NodeByClientKey(clientKey models.Key) *MMPoolNode {
	mmp.mu.Lock()
	defer mmp.mu.Unlock()
	node, _ := mmp.nodeByClientKey[clientKey]
	return node
}

func (mmp *MatchmakingPool) AddClient(client *models.ClientProfile, timeControl *models.TimeControl) error {
	mmp.mu.Lock()
	defer mmp.mu.Unlock()
	if _, ok := mmp.nodeByClientKey[client.ClientKey]; ok {
		return fmt.Errorf("client with key %s already in pool", client.ClientKey)
	}
	node := NewMMPoolNode(client, timeControl)
	if mmp.tail == nil {
		mmp.tail = node
		mmp.head = node
	} else {
		mmp.tail.next = node
		node.prev = mmp.tail
		mmp.tail = node
	}
	mmp.nodeByClientKey[client.ClientKey] = node
	return nil
}

func (mmp *MatchmakingPool) RemoveClient(clientKey models.Key) error {
	mmp.mu.Lock()
	defer mmp.mu.Unlock()
	node, ok := mmp.nodeByClientKey[clientKey]
	if !ok {
		return fmt.Errorf("client with key %s not in pool", clientKey)
	}
	if node.prev == nil {
		mmp.head = node.next
	} else {
		node.prev.next = node.next
	}
	if node.next == nil {
		mmp.tail = node.prev
	} else {
		node.next.prev = node.prev
	}
	delete(mmp.nodeByClientKey, clientKey)
	return nil
}

func (mmp *MatchmakingPool) GetBestMatch(node *MMPoolNode, waitTime int64) (*models.ClientProfile, error) {
	bestMatchPoolNode := node.next
	bestMatchWeight := float64(10000000000)
	nextPoolNode := node.next
	for nextPoolNode != nil {
		nextPoolNodeMatchWeight := weightProfileDiff(node.clientProfile, nextPoolNode.clientProfile, waitTime)

		if !IsMatchable(node.clientProfile, nextPoolNode.clientProfile, waitTime) {
			nextPoolNode = nextPoolNode.next
			continue
		}

		if nextPoolNodeMatchWeight < bestMatchWeight {
			bestMatchPoolNode = nextPoolNode
			bestMatchWeight = nextPoolNodeMatchWeight
		}
		nextPoolNode = nextPoolNode.next
	}

	if bestMatchPoolNode == nil {
		return nil, fmt.Errorf("no possible matches found")
	}
	return bestMatchPoolNode.clientProfile, nil
}

func weightProfileDiff(p1 *models.ClientProfile, p2 *models.ClientProfile, longestWaitSeconds int64) float64 {
	eloDiff := math.Abs(float64(p1.Elo - p2.Elo))
	eloCoeff := 100 / (float64(longestWaitSeconds) + 50) // asymptotic curve approaches 0, 2 @ t=0, 1 @ t=50, 0.5 @ t=100
	return eloDiff * eloCoeff
}

func IsMatchable(clientA *models.ClientProfile, clientB *models.ClientProfile, longestWaitSeconds int64) bool {
	return weightProfileDiff(clientA, clientB, longestWaitSeconds) <= 100
}
