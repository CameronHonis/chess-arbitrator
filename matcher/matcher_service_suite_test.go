package matcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var T *testing.T

func TestMatchService(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "MatcherService Suite")
}
