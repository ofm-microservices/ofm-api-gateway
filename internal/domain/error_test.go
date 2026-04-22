package gateway

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDomain(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Suite")
}

var _ = Describe("RegistrationConflictError", func() {
	It("returns a stable error string", func() {
		Expect((&RegistrationConflictError{}).Error()).To(Equal("registration conflict"))
	})
})
