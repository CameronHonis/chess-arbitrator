package clients_manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserClientsService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClientsManager Suite")
}
