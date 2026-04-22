package main

import (
	appfx "api-gateway/internal/fx"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

func TestMainPackage(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

type runnerStub struct {
	runCalls int
}

func (r *runnerStub) Run() {
	r.runCalls++
}

var _ = Describe("main", func() {
	It("builds the FX app with all gateway modules and runs it", func() {
		prevNewApp := newApp
		defer func() { newApp = prevNewApp }()

		stub := &runnerStub{}
		var captured []fx.Option
		newApp = func(opts ...fx.Option) runner {
			captured = opts
			return stub
		}

		main()

		Expect(stub.runCalls).To(Equal(1))
		Expect(captured).To(HaveLen(7))
		Expect(captured[0]).To(Equal(appfx.ConfigModule))
		Expect(captured[1]).To(Equal(appfx.LoggerModule))
		Expect(captured[2]).To(Equal(appfx.AppModule))
		Expect(captured[3]).To(Equal(appfx.MessagingModule))
		Expect(captured[4]).To(Equal(appfx.ServiceModule))
		Expect(captured[5]).To(Equal(appfx.HTTPModule))
		Expect(captured[6]).To(Equal(appfx.HTTPV1Module))
	})
})
