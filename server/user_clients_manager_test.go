package server

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserClientsManager", func() {
	var ucm *UserClientsManager
	var clientKey string
	BeforeEach(func() {
		ucm = NewUserClientsManager()
		clientKey = "some-client-key"
	})
	Describe("::AddClient", func() {
		It("adds the client key to the set", func() {
			err := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.clientKeys.Size()).To(Equal(1))
		})
		It("adds the channel keyed by the client", func() {
			err := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(err).ToNot(HaveOccurred())
			_, ok := ucm.channelByClientKey[clientKey]
			Expect(ok).To(BeTrue())
		})
		When("the client already exists", func() {
			It("returns an error", func() {
				err := ucm.AddClient(clientKey, make(chan *Prompt))
				Expect(err).ToNot(HaveOccurred())
				err2 := ucm.AddClient(clientKey, make(chan *Prompt))
				Expect(err2).To(HaveOccurred())
				Expect(err2).To(Equal(fmt.Errorf("client %s already exists", clientKey)))
			})
		})
	})
	Describe("::RemoveClient", func() {
		BeforeEach(func() {
			err := ucm.AddClient(clientKey, make(chan *Prompt))
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.clientKeys.Has(clientKey)).To(BeTrue())
			Expect(ucm.channelByClientKey[clientKey]).ToNot(BeNil())
		})
		It("removes the client key from the set", func() {
			err := ucm.RemoveClient(clientKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.clientKeys.Size()).To(Equal(0))
		})
		It("removes the channel keyed by the client", func() {
			err := ucm.RemoveClient(clientKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(ucm.channelByClientKey[clientKey]).To(BeNil())
		})
	})
})
