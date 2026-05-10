package gateway

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDomain(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gateway Domain Suite")
}

var _ = Describe("GigConflictError", func() {
	It("renders a stable error string", func() {
		Expect((&GigConflictError{State: "draft"}).Error()).To(Equal("gig conflict"))
	})
})
