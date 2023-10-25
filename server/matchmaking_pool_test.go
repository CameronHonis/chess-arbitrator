package server

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchmakingPool", func() {
	var matchmakingPool *MatchmakingPool
	BeforeEach(func() {
		matchmakingPool = NewMatchmakingPool()
	})
	Describe("AddClient", func() {
		var clientProfile *ClientProfile
		BeforeEach(func() {
			clientProfile = NewClientProfile("some-client-key", 1000)
		})
		Context("when the client is not already in the pool", func() {
			It("should add the client to the pool", func() {
				err := matchmakingPool.AddClient(clientProfile)
				Expect(err).To(BeNil())
				Expect(matchmakingPool.head).ToNot(BeNil())
				Expect(matchmakingPool.head).To(Equal(matchmakingPool.tail))
				node, ok := matchmakingPool.nodeByClientKey[clientProfile.ClientKey]
				Expect(ok).To(BeTrue())
				Expect(node).To(Equal(matchmakingPool.head))
			})
		})
		Context("when the client is already in the pool", func() {
			BeforeEach(func() {
				mMPoolNode := &MMPoolNode{
					clientProfile: clientProfile,
				}
				matchmakingPool.nodeByClientKey[clientProfile.ClientKey] = mMPoolNode
				matchmakingPool.head = mMPoolNode
				matchmakingPool.tail = mMPoolNode
			})
			It("should return an error", func() {
				err := matchmakingPool.AddClient(clientProfile)
				Expect(err).To(Equal(fmt.Errorf("client with key %s already in pool", clientProfile.ClientKey)))
			})
		})
		Context("when a client is already in the pool", func() {
			var otherClientProfile *ClientProfile
			BeforeEach(func() {
				otherClientProfile = NewClientProfile("some-other-client-key", 1000)
				mMPoolNode := &MMPoolNode{
					clientProfile: otherClientProfile,
				}
				matchmakingPool.nodeByClientKey[otherClientProfile.ClientKey] = mMPoolNode
				matchmakingPool.head = mMPoolNode
				matchmakingPool.tail = mMPoolNode
				Expect(matchmakingPool.head).ToNot(BeNil())
				Expect(matchmakingPool.head).To(Equal(matchmakingPool.tail))
				Expect(matchmakingPool.head.clientProfile).To(Equal(otherClientProfile))
			})
			It("should add the client to the pool", func() {
				err := matchmakingPool.AddClient(clientProfile)
				Expect(err).To(BeNil())
				Expect(matchmakingPool.head.clientProfile).To(Equal(otherClientProfile))
				Expect(matchmakingPool.head.next).ToNot(BeNil())
				Expect(matchmakingPool.head.next.clientProfile).To(Equal(clientProfile))
				Expect(matchmakingPool.head.next.next).To(BeNil())
				Expect(matchmakingPool.tail).To(Equal(matchmakingPool.head.next))
				Expect(matchmakingPool.tail.prev).To(Equal(matchmakingPool.head))
				node, ok := matchmakingPool.nodeByClientKey[clientProfile.ClientKey]
				Expect(ok).To(BeTrue())
				Expect(node).To(Equal(matchmakingPool.tail))
			})
		})
	})
	Describe("RemoveClient", func() {
		var clientA, clientB, clientC *ClientProfile
		BeforeEach(func() {
			clientA = NewClientProfile("client-key-a", 1000)
			clientB = NewClientProfile("client-key-b", 1000)
			clientC = NewClientProfile("client-key-c", 1000)
			err := matchmakingPool.AddClient(clientA)
			Expect(err).To(BeNil())
			err = matchmakingPool.AddClient(clientB)
			Expect(err).To(BeNil())
			err = matchmakingPool.AddClient(clientC)
			Expect(err).To(BeNil())
			Expect(matchmakingPool.head).ToNot(BeNil())
			Expect(matchmakingPool.head.clientProfile).To(Equal(clientA))
			Expect(matchmakingPool.tail).ToNot(BeNil())
			Expect(matchmakingPool.tail.clientProfile).To(Equal(clientC))
			Expect(matchmakingPool.head.next).ToNot(BeNil())
			Expect(matchmakingPool.head.next.clientProfile).To(Equal(clientB))
			Expect(matchmakingPool.tail.prev).ToNot(BeNil())
			Expect(matchmakingPool.tail.prev.clientProfile).To(Equal(clientB))
			Expect(matchmakingPool.head.next.next).ToNot(BeNil())
			Expect(matchmakingPool.head.next.next.clientProfile).To(Equal(clientC))
			Expect(matchmakingPool.tail.prev.prev).ToNot(BeNil())
			Expect(matchmakingPool.tail.prev.prev.clientProfile).To(Equal(clientA))
		})
		Context("when the client is the head of the pool", func() {
			It("removes the client and re-assign the head", func() {
				err := matchmakingPool.RemoveClient(clientA.ClientKey)
				Expect(err).To(BeNil())
				Expect(matchmakingPool.head).ToNot(BeNil())
				Expect(matchmakingPool.head.clientProfile).To(Equal(clientB))
				Expect(matchmakingPool.head.prev).To(BeNil())
				_, ok := matchmakingPool.nodeByClientKey[clientA.ClientKey]
				Expect(ok).To(BeFalse())
			})
		})
		Context("when the client is the tail of the pool", func() {
			It("removes the client and re-assign the tail", func() {
				err := matchmakingPool.RemoveClient(clientC.ClientKey)
				Expect(err).To(BeNil())
				Expect(matchmakingPool.tail).ToNot(BeNil())
				Expect(matchmakingPool.tail.clientProfile).To(Equal(clientB))
				Expect(matchmakingPool.tail.next).To(BeNil())
				_, ok := matchmakingPool.nodeByClientKey[clientC.ClientKey]
				Expect(ok).To(BeFalse())
			})
		})
		Context("when the client is in the middle of the pool", func() {
			It("removes the client and modifies the pointers of its neighbors", func() {
				err := matchmakingPool.RemoveClient(clientB.ClientKey)
				Expect(err).To(BeNil())
				Expect(matchmakingPool.head).ToNot(BeNil())
				Expect(matchmakingPool.head.clientProfile).To(Equal(clientA))
				Expect(matchmakingPool.tail).ToNot(BeNil())
				Expect(matchmakingPool.tail.clientProfile).To(Equal(clientC))
				Expect(matchmakingPool.head.next).ToNot(BeNil())
				Expect(matchmakingPool.head.next.clientProfile).To(Equal(clientC))
				Expect(matchmakingPool.tail.prev).ToNot(BeNil())
				Expect(matchmakingPool.tail.prev.clientProfile).To(Equal(clientA))
				Expect(matchmakingPool.head.prev).To(BeNil())
				Expect(matchmakingPool.tail.next).To(BeNil())
				_, ok := matchmakingPool.nodeByClientKey[clientB.ClientKey]
				Expect(ok).To(BeFalse())
			})
		})
		Context("when the client is not in the pool", func() {
			It("returns an error", func() {
				err := matchmakingPool.RemoveClient("some-non-existent-client-key")
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(fmt.Errorf("client with key some-non-existent-client-key not in pool")))
			})
		})
	})

})
