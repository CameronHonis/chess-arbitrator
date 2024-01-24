package helpers_test

import (
	"github.com/CameronHonis/chess-arbitrator/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Introspect", func() {
	Describe("IsClientKey", func() {
		When("the input is a client key", func() {
			It("returns true", func() {
				isKey := helpers.IsClientKey("4efdac2c21854bd57b72f6798d0fe646fecb9c36c4194da410e315027a8c0847")
				Expect(isKey).To(BeTrue())
			})
		})
		When("the input is not length 64", func() {
			It("returns false", func() {
				isKey := helpers.IsClientKey("4EFDAC2C21854BD57B72F6798D0FE646FECB9C36C4194DA")
				Expect(isKey).To(BeFalse())
			})
		})
		When("the input has invalid key chars", func() {
			It("returns false", func() {
				isKey := helpers.IsClientKey("**FDAC2C21854BD57B72F6798D0FE646FECB9C36C4194DA410E315027A8C0847")
				Expect(isKey).To(BeFalse())
			})
		})
	})
})
