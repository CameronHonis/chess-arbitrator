package matchmaking_test

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MatchmakingPool", func() {
	var matchmakingPool *matchmaking.MatchmakingPool
	BeforeEach(func() {
		matchmakingPool = matchmaking.NewMatchmakingPool()
	})
	Describe("AddClient", func() {
		var clientProfile *models.ClientProfile
		BeforeEach(func() {
			clientProfile = models.NewClientProfile("some-client-key", 1000)
		})
		Context("when the client is not already in the pool", func() {
			It("should add the client to the pool", func() {
				err := matchmakingPool.AddClient(clientProfile, builders.NewBlitzTimeControl())
				Expect(err).To(BeNil())
				Expect(matchmakingPool.Head()).To(Equal(matchmakingPool.Tail()))
				Expect(matchmakingPool.NodeByClientKey(clientProfile.ClientKey)).To(Equal(matchmakingPool.Head()))
			})
		})
		Context("when the client is already in the pool", func() {
			BeforeEach(func() {
				Expect(matchmakingPool.AddClient(clientProfile, builders.NewBlitzTimeControl())).ToNot(HaveOccurred())
			})
			It("should return an error", func() {
				err := matchmakingPool.AddClient(clientProfile, builders.NewBlitzTimeControl())
				Expect(err).To(Equal(fmt.Errorf("client with key %s already in pool", clientProfile.ClientKey)))
			})
		})
		Context("when a client is already in the pool", func() {
			var otherClientProfile *models.ClientProfile
			BeforeEach(func() {
				otherClientProfile = models.NewClientProfile("some-other-client-key", 1000)
				Expect(matchmakingPool.AddClient(otherClientProfile, builders.NewBlitzTimeControl())).ToNot(HaveOccurred())
			})
			It("should add the client to the pool", func() {
				Expect(matchmakingPool.AddClient(clientProfile, builders.NewBlitzTimeControl())).ToNot(HaveOccurred())
				Expect(matchmakingPool.Head().ClientProfile()).To(Equal(otherClientProfile))
				Expect(matchmakingPool.Head().Next().ClientProfile()).To(Equal(clientProfile))
				Expect(matchmakingPool.Tail().ClientProfile()).To(Equal(clientProfile))
				Expect(matchmakingPool.NodeByClientKey(otherClientProfile.ClientKey).ClientProfile()).To(Equal(otherClientProfile))
				Expect(matchmakingPool.NodeByClientKey(clientProfile.ClientKey).ClientProfile()).To(Equal(clientProfile))
			})
		})
	})
	Describe("RemoveClient", func() {
		var clientA, clientB, clientC *models.ClientProfile
		BeforeEach(func() {
			clientA = models.NewClientProfile("client-key-a", 1000)
			clientB = models.NewClientProfile("client-key-b", 1000)
			clientC = models.NewClientProfile("client-key-c", 1000)
			Expect(matchmakingPool.AddClient(clientA, builders.NewBlitzTimeControl())).ToNot(HaveOccurred())
			Expect(matchmakingPool.AddClient(clientB, builders.NewBlitzTimeControl())).ToNot(HaveOccurred())
			Expect(matchmakingPool.AddClient(clientC, builders.NewBlitzTimeControl())).ToNot(HaveOccurred())
		})
		Context("when the client is the head of the pool", func() {
			It("removes the client and re-assign the head", func() {
				Expect(matchmakingPool.RemoveClient(clientA.ClientKey)).ToNot(HaveOccurred())
				Expect(matchmakingPool.Head().ClientProfile()).To(Equal(clientB))
				Expect(matchmakingPool.Head().Prev()).To(BeNil())
				Expect(matchmakingPool.NodeByClientKey(clientA.ClientKey))
			})
		})
		Context("when the client is the tail of the pool", func() {
			It("removes the client and re-assign the tail", func() {
				Expect(matchmakingPool.RemoveClient(clientC.ClientKey)).ToNot(HaveOccurred())
				Expect(matchmakingPool.Tail().ClientProfile()).To(Equal(clientB))
				Expect(matchmakingPool.Tail().Next()).To(BeNil())
				Expect(matchmakingPool.NodeByClientKey(clientC.ClientKey)).To(BeNil())
			})
		})
		Context("when the client is in the middle of the pool", func() {
			It("removes the client and modifies the pointers of its neighbors", func() {
				Expect(matchmakingPool.RemoveClient(clientB.ClientKey)).ToNot(HaveOccurred())
				Expect(matchmakingPool.Head().ClientProfile()).To(Equal(clientA))
				Expect(matchmakingPool.Tail().ClientProfile()).To(Equal(clientC))
				Expect(matchmakingPool.Head().Next().ClientProfile()).To(Equal(clientC))
				Expect(matchmakingPool.Tail().Prev().ClientProfile()).To(Equal(clientA))
				Expect(matchmakingPool.NodeByClientKey(clientB.ClientKey)).To(BeNil())
			})
		})
		Context("when the client is not in the pool", func() {
			It("returns an error", func() {
				Expect(matchmakingPool.RemoveClient("some-non-existent-client-key")).To(HaveOccurred())
			})
		})
	})
})
