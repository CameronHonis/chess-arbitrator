package helpers_test

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/helpers"
	"github.com/CameronHonis/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("LogWrappers", func() {
	Describe("PrettyClientDecoratorRule", func() {
		var rule log.DecoratorRule
		When("the env is indicative of a client key", func() {
			var isClientKeyArg string
			BeforeEach(func() {
				rule = helpers.PrettyClientDecoratorRule(func(s string) bool {
					isClientKeyArg = s
					return true
				})
			})
			It("calls IsClientKey", func() {
				_ = rule("ABCD123XYZ")
				Expect(isClientKeyArg).To(Equal("abcd123xyz"))
			})
			It("returns a shortened version of the env tag", func() {
				decorator := rule("ABCD123XYZ")
				decoratorAddr := reflect.ValueOf(decorator).Pointer()
				expAddr := reflect.ValueOf(helpers.PrettyClientDecorator).Pointer()
				Expect(decoratorAddr).To(Equal(expAddr))
			})
		})
		When("the env is not indicative of a client key", func() {
			BeforeEach(func() {
				rule = helpers.PrettyClientDecoratorRule(func(_ string) bool {
					return false
				})
			})
			It("returns nil", func() {
				decorator := rule("ASDF")
				Expect(decorator).To(BeNil())
			})
		})
	})
	Describe("PrettyClientDecorator", func() {
		It("returns a minified client key", func() {
			expEnvName := fmt.Sprintf("%s[abc..xyz]%s", "\x1b[36m", "\x1b[0m")
			Expect(helpers.PrettyClientDecorator("[abcdef111111111111111111tuvwxyz]")).To(Equal(expEnvName))
		})
	})
})
