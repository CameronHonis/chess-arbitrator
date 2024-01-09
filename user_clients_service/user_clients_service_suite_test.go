package user_clients_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserClientsService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserClientsService Suite")
}
