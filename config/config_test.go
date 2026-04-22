package config

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Load", func() {
	var prevWD string

	BeforeEach(func() {
		var err error
		prevWD, err = os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		tmpDir := GinkgoT().TempDir()
		Expect(os.Chdir(tmpDir)).To(Succeed())

		for _, key := range []string{
			"APP_ENV",
			"LOG_LEVEL",
			"HTTP_HOST",
			"HTTP_PORT",
			"REGISTRATION_SAGA_ADDRESS",
		} {
			Expect(os.Unsetenv(key)).To(Succeed())
		}
	})

	AfterEach(func() {
		Expect(os.Chdir(prevWD)).To(Succeed())
	})

	It("loads defaults when no environment overrides are present", func() {
		cfg, err := Load()

		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.App.Env).To(Equal("local"))
		Expect(cfg.App.LogLevel).To(Equal("info"))
		Expect(cfg.HTTP.Host).To(Equal("0.0.0.0"))
		Expect(cfg.HTTP.Port).To(Equal(8080))
		Expect(cfg.RegistrationSaga.Address).To(Equal("127.0.0.1:9090"))
	})

	It("loads explicit environment overrides", func() {
		Expect(os.Setenv("APP_ENV", "test")).To(Succeed())
		Expect(os.Setenv("LOG_LEVEL", "debug")).To(Succeed())
		Expect(os.Setenv("HTTP_HOST", "127.0.0.1")).To(Succeed())
		Expect(os.Setenv("HTTP_PORT", "9091")).To(Succeed())
		Expect(os.Setenv("REGISTRATION_SAGA_ADDRESS", "127.0.0.1:9191")).To(Succeed())

		cfg, err := Load()

		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.App.Env).To(Equal("test"))
		Expect(cfg.App.LogLevel).To(Equal("debug"))
		Expect(cfg.HTTP.Host).To(Equal("127.0.0.1"))
		Expect(cfg.HTTP.Port).To(Equal(9091))
		Expect(cfg.RegistrationSaga.Address).To(Equal("127.0.0.1:9191"))
	})

	It("wraps env parsing failures", func() {
		Expect(os.Setenv("HTTP_PORT", "bad-port")).To(Succeed())

		cfg, err := Load()

		Expect(cfg).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("parse env config"))
	})
})
