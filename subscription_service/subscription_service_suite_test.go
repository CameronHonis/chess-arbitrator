package subscription_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSubscriptionService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SubscriptionService Suite")
}
