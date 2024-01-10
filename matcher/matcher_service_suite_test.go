package matcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMatchService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MatcherService Suite")
}
