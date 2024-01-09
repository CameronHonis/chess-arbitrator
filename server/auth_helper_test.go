package server_test

import (
	. "github.com/CameronHonis/chess-arbitrator/helpers"
	. "github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	Describe("GenerateKeyset", func() {
		It("generates a private key and a public key", func() {
			publicKey, privateKey := GenerateKeyset()
			Expect(publicKey).ToNot(BeEmpty())
			Expect(privateKey).ToNot(BeEmpty())
			Expect(len(privateKey)).To(Equal(36), "private key should be a uuid")
			Expect(len(publicKey)).To(Equal(64), "public key should be a sha256 hex hash")
		})
		It("generates a unique public & private key each time", func() {
			pubKeySet := set.EmptySet[Key]()
			privKeySet := set.EmptySet[Key]()
			for i := 0; i < 1000; i++ {
				publicKey, privateKey := GenerateKeyset()
				Expect(pubKeySet.Has(publicKey)).To(BeFalse())
				Expect(privKeySet.Has(privateKey)).To(BeFalse())
				pubKeySet.Add(publicKey)
				privKeySet.Add(privateKey)
			}
		})
		It("generates a public key that is a hex encoded sha256 hash of the private key", func() {
			publicKey, privateKey := GenerateKeyset()
			Expect(ValidatePrivateKey(publicKey, privateKey)).To(BeTrue())
		})
	})
})
