package chess_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestChess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chess Suite")
}
