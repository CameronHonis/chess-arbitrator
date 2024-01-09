package matchmaking

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
	"sync"
	"time"
)

type MMPoolNode struct {
	next          *MMPoolNode
	prev          *MMPoolNode
	clientProfile *models.ClientProfile
	timeJoined    int64
}

func NewMMPoolNode(profile *models.ClientProfile) *MMPoolNode {
	return &MMPoolNode{
		clientProfile: profile,
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
	// doubly linked list of client profiles, sorted by queue time
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

func (mmp *MatchmakingPool) AddClient(client *models.ClientProfile) error {
	mmp.mu.Lock()
	if _, ok := mmp.nodeByClientKey[client.ClientKey]; ok {
		return fmt.Errorf("client with key %s already in pool", client.ClientKey)
	}
	node := NewMMPoolNode(client)
	if mmp.tail == nil {
		mmp.tail = node
		mmp.head = node
	} else {
		mmp.tail.next = node
		node.prev = mmp.tail
		mmp.tail = node
	}
	mmp.nodeByClientKey[client.ClientKey] = node
	mmp.mu.Unlock()
	return nil
}

func (mmp *MatchmakingPool) RemoveClient(clientKey models.Key) error {
	mmp.mu.Lock()
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
	mmp.mu.Unlock()
	return nil
}
