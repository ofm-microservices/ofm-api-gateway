package appfx

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	grpcclient "api-gateway/internal/presentation/grpc"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	registrationv1 "github.com/ofm-microseervices/ofm-common/proto/registration/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func TestFX(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "FX Suite")
}

type lifecycleStub struct {
	hooks []fx.Hook
}

func (l *lifecycleStub) Append(h fx.Hook) {
	l.hooks = append(l.hooks, h)
}

type registrationPublisherStub struct{}

func (registrationPublisherStub) StartRegistration(context.Context, gateway.SignUpRequest) (*gateway.SignUpResult, error) {
	return &gateway.SignUpResult{Status: "pending"}, nil
}

type registrationServiceStub struct{}

func (registrationServiceStub) SignUp(context.Context, gateway.SignUpRequest) (*gateway.SignUpResult, error) {
	return &gateway.SignUpResult{Status: "pending"}, nil
}

type authHandlerStub struct {
	registered bool
}

func (h *authHandlerStub) RegisterRoutes(router fiber.Router) {
	h.registered = true
	router.Post("/auth/sign-up", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})
}

func (h *authHandlerStub) HandleSignUp(*fiber.Ctx) error {
	return nil
}

type httpServerStub struct {
	app           *fiber.App
	startCalls    int
	shutdownCalls int
}

func (s *httpServerStub) Start() error {
	s.startCalls++
	return nil
}

func (s *httpServerStub) Shutdown(context.Context) error {
	s.shutdownCalls++
	return nil
}

func (s *httpServerStub) App() *fiber.App {
	return s.app
}

type fakeRegistrationServer struct {
	registrationv1.UnimplementedRegistrationServiceServer
}

func (fakeRegistrationServer) StartRegistration(context.Context, *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error) {
	return &registrationv1.StartRegistrationResponse{Status: "pending"}, nil
}

var _ = Describe("FX providers", func() {
	var (
		cfg *config.Config
		lg  logging.Logger
	)

	BeforeEach(func() {
		var err error
		lg, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())

		cfg = &config.Config{
			App:              config.AppConfig{Env: "test", LogLevel: "debug"},
			HTTP:             config.HTTPConfig{Host: "127.0.0.1", Port: 8080},
			RegistrationSaga: config.RegistrationSagaConfig{Address: "127.0.0.1:9090"},
		}
	})

	It("loads config from the current environment", func() {
		prevWD, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			Expect(os.Chdir(prevWD)).To(Succeed())
			Expect(os.Unsetenv("HTTP_PORT")).To(Succeed())
		}()

		tmpDir := GinkgoT().TempDir()
		Expect(os.Chdir(tmpDir)).To(Succeed())
		Expect(os.Setenv("HTTP_PORT", "9099")).To(Succeed())

		loaded, err := ProvideConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(loaded.HTTP.Port).To(Equal(9099))
	})

	It("constructs a logger and registers an OnStop hook", func() {
		lc := &lifecycleStub{}

		logger, err := ProvideLogger(lc, cfg)

		Expect(err).NotTo(HaveOccurred())
		Expect(logger).NotTo(BeNil())
		Expect(lc.hooks).To(HaveLen(1))
		Expect(lc.hooks[0].OnStop).NotTo(BeNil())
		Expect(func() {
			_ = lc.hooks[0].OnStop(context.Background())
		}).NotTo(Panic())
	})

	It("returns a logger error for an invalid log level", func() {
		lc := &lifecycleStub{}
		cfg.App.LogLevel = "bad-level"

		logger, err := ProvideLogger(lc, cfg)

		Expect(logger).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(lc.hooks).To(BeEmpty())
	})

	It("constructs the registration application service", func() {
		svc, err := ProvideRegistrationService(registrationPublisherStub{}, lg)

		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())
	})

	It("constructs the HTTP server", func() {
		httpSrv, err := ProvideHTTPServer(cfg, lg)
		Expect(err).NotTo(HaveOccurred())
		Expect(httpSrv).NotTo(BeNil())
	})

	It("constructs the v1 auth handler", func() {
		handler, err := ProvideHTTPV1AuthHandler(registrationServiceStub{}, lg)

		Expect(err).NotTo(HaveOccurred())
		Expect(handler).NotTo(BeNil())
	})

	It("registers versioned routes", func() {
		srv := &httpServerStub{app: fiber.New()}
		handler := &authHandlerStub{}
		InvokeRegisterHTTPV1Routes(srv, handler)

		Expect(handler.registered).To(BeTrue())
	})

	It("wires HTTP server lifecycle hooks", func() {
		lc := &lifecycleStub{}
		srv := &httpServerStub{app: fiber.New()}

		InvokeRunHTTPServer(lc, srv, lg)

		Expect(lc.hooks).To(HaveLen(1))
		Expect(lc.hooks[0].OnStart(context.Background())).To(Succeed())
		Eventually(func() int {
			return srv.startCalls
		}, time.Second, 20*time.Millisecond).Should(Equal(1))
		Expect(lc.hooks[0].OnStop(context.Background())).To(Succeed())
		Expect(srv.shutdownCalls).To(Equal(1))
	})

	It("constructs a grpc registration publisher and registers its shutdown hook", func() {
		server, address := startRegistrationServer()
		defer server.Stop()

		lc := &lifecycleStub{}
		cfg.RegistrationSaga.Address = address

		pub, err := ProvideRegistrationPublisher(lc, cfg, lg)

		Expect(err).NotTo(HaveOccurred())
		Expect(pub).NotTo(BeNil())
		Expect(lc.hooks).To(HaveLen(1))
		Expect(lc.hooks[0].OnStop(context.Background())).To(Succeed())
	})

	It("returns grpc dial errors from registration publisher construction", func() {
		lc := &lifecycleStub{}
		cfg.RegistrationSaga.Address = ""

		pub, err := ProvideRegistrationPublisher(lc, cfg, lg)

		Expect(pub).To(BeNil())
		Expect(err).To(MatchError(grpcclient.ErrEmptyAddress))
		Expect(lc.hooks).To(BeEmpty())
	})

	It("emits the start log without panicking", func() {
		Expect(func() { InvokeStartLog(lg) }).NotTo(Panic())
	})
})

func startRegistrationServer() (*grpc.Server, string) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	Expect(err).NotTo(HaveOccurred())

	server := grpc.NewServer()
	registrationv1.RegisterRegistrationServiceServer(server, fakeRegistrationServer{})

	go func() {
		_ = server.Serve(lis)
	}()

	return server, lis.Addr().String()
}
