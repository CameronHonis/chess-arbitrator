package clients_manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var T *testing.T

func TestClientsManager(t *testing.T) {
	T = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClientsManager Suite")
}
