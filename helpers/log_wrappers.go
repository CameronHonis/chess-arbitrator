package helpers

import (
	"fmt"
	"github.com/CameronHonis/log"
	"strings"
)

func PrettyClientDecoratorRule(isClientKey func(string) bool) log.DecoratorRule {
	return func(envName string) log.Decorator {
		if isClientKey(strings.ToLower(envName)) {
			return PrettyClientDecorator
		}
		return nil
	}
}

func PrettyClientDecorator(envName string) string {
	return log.WrapCyan(fmt.Sprintf("%s..%s", envName[:4], envName[len(envName)-4:]))
}
