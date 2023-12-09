package server

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserClientsManager", func() {
	var client *UserClient
	var clientKey string
	BeforeEach(func() {
		subscriptionManager = nil
		userClientsManager = nil
		GetUserClientsManager()
		client = NewUserClient(nil, func(client *UserClient) {})
		client.publicKey = "some-public-key"
		clientKey = client.publicKey
	})
	Describe("::AddClient", func() {
		When("the client hasn't been added", func() {
			It("adds the client to the map", func() {
				err := userClientsManager.AddClient(client)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(userClientsManager.clientByPublicKey)).To(Equal(1))
				client, ok := userClientsManager.clientByPublicKey[clientKey]
				Expect(ok).To(BeTrue())
				Expect(*client).To(BeAssignableToTypeOf(UserClient{}))
				Expect(client.publicKey).To(Equal(clientKey))
				Expect(client.inChannel).ToNot(BeNil())
				Expect(client.outChannel).ToNot(BeNil())
			})
		})
		When("the client already exists", func() {
			BeforeEach(func() {
				err := userClientsManager.AddClient(client)
				Expect(err).ToNot(HaveOccurred())
			})
			It("returns an error", func() {
				err2 := userClientsManager.AddClient(client)
				Expect(err2).To(HaveOccurred())
				Expect(err2).To(Equal(fmt.Errorf("client with key %s already exists", clientKey)))
			})
		})
	})
	Describe("::RemoveClient", func() {
		var topic MessageTopic
		When("a player was added", func() {
			BeforeEach(func() {
				topic = "auth"
				addClientErr := userClientsManager.AddClient(client)
				Expect(addClientErr).ToNot(HaveOccurred())
				_, ok := userClientsManager.clientByPublicKey[clientKey]
				Expect(ok).To(BeTrue())

				subClientErr := subscriptionManager.SubClientTo(clientKey, topic)
				Expect(subClientErr).ToNot(HaveOccurred())
				Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(1))
				Expect(subscriptionManager.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(1))
			})
			It("removes the client from the client map", func() {
				err := userClientsManager.RemoveClient(client)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(userClientsManager.clientByPublicKey)).To(Equal(0))
			})
			It("unsubs client from all topics", func() {
				err := userClientsManager.RemoveClient(client)
				Expect(err).ToNot(HaveOccurred())
				Expect(subscriptionManager.GetSubbedTopics(clientKey).Size()).To(Equal(0))
				Expect(subscriptionManager.GetClientKeysSubbedToTopic(topic).Size()).To(Equal(0))
			})
		})
		When("the player was never added", func() {
			It("returns an error", func() {
				err := userClientsManager.RemoveClient(client)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("client with key %s isn't an established client", clientKey)))
			})
		})
	})
	Describe("::GetAllOutChannels", func() {
		// is this flakey?
		var otherClient *UserClient
		BeforeEach(func() {
			addClientErr := userClientsManager.AddClient(client)
			Expect(addClientErr).ToNot(HaveOccurred())
			otherClient = NewUserClient(nil, func(client *UserClient) {})
			otherClient.publicKey = "other-public-key"
			addClientErr = userClientsManager.AddClient(otherClient)
			Expect(addClientErr).ToNot(HaveOccurred())
		})
		It("returns a slice of all client channels", func() {
			channels := userClientsManager.GetAllOutChannels()
			Expect(channels).To(HaveLen(2))
			Expect(channels[client.PublicKey()]).To(Equal(client.outChannel))
			Expect(channels[otherClient.PublicKey()]).To(Equal(otherClient.outChannel))
		})
	})
})
