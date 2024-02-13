package matchmaking_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var T *testing.T

func TestMatchmaking(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Matchmaking Suite")
}
