package server

import (
	"fmt"
	"time"
)

type MMPoolNode struct {
	next          *MMPoolNode
	prev          *MMPoolNode
	clientProfile *ClientProfile
	timeJoined    int64
}

type MatchmakingPool struct {
	// doubly linked list of client profiles, sorted by queue time
	// prioritizes add/remove speed
	head *MMPoolNode
	tail *MMPoolNode
	// map to allow for O(1) lookup time of nodes by client key
	nodeByClientKey map[Key]*MMPoolNode
}

func NewMatchmakingPool() *MatchmakingPool {
	return &MatchmakingPool{
		nodeByClientKey: make(map[Key]*MMPoolNode),
	}
}

func (mmp *MatchmakingPool) AddClient(client *ClientProfile) error {
	if _, ok := mmp.nodeByClientKey[client.ClientKey]; ok {
		return fmt.Errorf("client with key %s already in pool", client.ClientKey)
	}
	node := &MMPoolNode{
		clientProfile: client,
		timeJoined:    time.Now().Unix(),
	}
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

func (mmp *MatchmakingPool) RemoveClient(clientKey Key) error {
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
