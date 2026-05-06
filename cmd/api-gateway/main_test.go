package main

import (
	"testing"

	"go.uber.org/fx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMainCmd(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

type runnerStub struct{ ran bool }

func (r *runnerStub) Run() { r.ran = true }

var _ = Describe("main wiring", func() {
	It("constructs and runs the fx app", func() {
		prev := newApp
		defer func() { newApp = prev }()

		rs := &runnerStub{}
		newApp = func(opts ...fx.Option) runner { return rs }

		main()
		Expect(rs.ran).To(BeTrue())
	})
})
