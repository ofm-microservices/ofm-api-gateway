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
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	authv1 "github.com/ofm-microservices/ofm-common/proto/auth/v1"
	registrationv1 "github.com/ofm-microservices/ofm-common/proto/registration/v1"
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

func (registrationPublisherStub) VerifyEmail(context.Context, gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error) {
	return &gateway.VerifyEmailResult{Status: "verifying_email"}, nil
}

func (registrationPublisherStub) GetRegistrationStatus(context.Context, string, string) (*gateway.RegistrationStatus, error) {
	return &gateway.RegistrationStatus{Status: "completed", UserID: "user-1"}, nil
}

type tokenIssuerStub struct{}

func (tokenIssuerStub) IssueRegistrationTokens(context.Context, string) (*gateway.CompleteRegistrationResult, error) {
	return &gateway.CompleteRegistrationResult{TokenType: "Bearer"}, nil
}

func (tokenIssuerStub) Close() error { return nil }

type authSessionServiceStub struct{}

func (authSessionServiceStub) SignIn(context.Context, gateway.SignInRequest) (*gateway.AuthTokensResult, error) {
	return &gateway.AuthTokensResult{TokenType: "Bearer"}, nil
}

func (authSessionServiceStub) Close() error { return nil }

type registrationServiceStub struct{}

func (registrationServiceStub) SignUp(context.Context, gateway.SignUpRequest) (*gateway.SignUpResult, error) {
	return &gateway.SignUpResult{Status: "pending"}, nil
}

func (registrationServiceStub) VerifyEmail(context.Context, gateway.VerifyEmailRequest) (*gateway.VerifyEmailResult, error) {
	return &gateway.VerifyEmailResult{Status: "verifying_email"}, nil
}

func (registrationServiceStub) CompleteRegistration(context.Context, gateway.CompleteRegistrationRequest) (*gateway.CompleteRegistrationResult, error) {
	return &gateway.CompleteRegistrationResult{TokenType: "Bearer"}, nil
}

type authHandlerStub struct {
	registered bool
}

func (h *authHandlerStub) RegisterRoutes(router fiber.Router) {
	h.registered = true
	router.Post("/auth/sign-up", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})
	router.Post("/auth/sign-in", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

func (h *authHandlerStub) HandleSignUp(*fiber.Ctx) error {
	return nil
}

func (h *authHandlerStub) HandleSignIn(*fiber.Ctx) error {
	return nil
}

func (h *authHandlerStub) HandleVerifyEmail(*fiber.Ctx) error {
	return nil
}

func (h *authHandlerStub) HandleCompleteRegistration(*fiber.Ctx) error {
	return nil
}

type gigHandlerStub struct {
	registered bool
}

func (h *gigHandlerStub) RegisterRoutes(router fiber.Router) {
	h.registered = true
	router.Post("/gigs/drafts", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})
}

func (h *gigHandlerStub) HandleCreateDraft(*fiber.Ctx) error      { return nil }
func (h *gigHandlerStub) HandleUpdateBasicInfo(*fiber.Ctx) error  { return nil }
func (h *gigHandlerStub) HandleReplacePackages(*fiber.Ctx) error  { return nil }
func (h *gigHandlerStub) HandleReplaceQuestions(*fiber.Ctx) error { return nil }
func (h *gigHandlerStub) HandleReplaceMedia(*fiber.Ctx) error     { return nil }
func (h *gigHandlerStub) HandleGetDraft(*fiber.Ctx) error         { return nil }
func (h *gigHandlerStub) HandleGetBySlug(*fiber.Ctx) error        { return nil }
func (h *gigHandlerStub) HandlePublish(*fiber.Ctx) error          { return nil }

type orderHandlerStub struct {
	registered bool
}

func (h *orderHandlerStub) RegisterRoutes(router fiber.Router) {
	h.registered = true
	router.Post("/orders", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})
}

func (h *orderHandlerStub) HandleCreateOrder(*fiber.Ctx) error     { return nil }
func (h *orderHandlerStub) HandleConfirmOrder(*fiber.Ctx) error    { return nil }
func (h *orderHandlerStub) HandleDeliverOrder(*fiber.Ctx) error    { return nil }
func (h *orderHandlerStub) HandleAcceptDelivery(*fiber.Ctx) error  { return nil }
func (h *orderHandlerStub) HandleRequestRevision(*fiber.Ctx) error { return nil }
func (h *orderHandlerStub) HandleOpenDispute(*fiber.Ctx) error     { return nil }

type reviewHandlerStub struct {
	registered bool
}

func (h *reviewHandlerStub) RegisterRoutes(router fiber.Router) {
	h.registered = true
	router.Post("/orders/:order_id/reviews", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusCreated)
	})
	router.Get("/orders/:order_id/review", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	router.Get("/gigs/:gig_id/reviews", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	router.Get("/users/:user_id/reviews", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

func (h *reviewHandlerStub) HandleCreateReview(*fiber.Ctx) error         { return nil }
func (h *reviewHandlerStub) HandleGetReviewByOrderID(*fiber.Ctx) error   { return nil }
func (h *reviewHandlerStub) HandleListReviewsByGigID(*fiber.Ctx) error   { return nil }
func (h *reviewHandlerStub) HandleListReviewsByBuyerID(*fiber.Ctx) error { return nil }

type searchHandlerStub struct {
	registered bool
}

func (h *searchHandlerStub) RegisterRoutes(router fiber.Router) {
	h.registered = true
	router.Get("/search", h.HandleSearch)
}

func (h *searchHandlerStub) HandleSearch(*fiber.Ctx) error { return nil }

type onboardingHandlerStub struct {
	registered bool
}

func (h *onboardingHandlerStub) RegisterRoutes(router fiber.Router) {
	h.registered = true
	router.Post("/freelancer/onboarding/start", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusAccepted)
	})
}

func (h *onboardingHandlerStub) HandleStartFreelancerOnboarding(*fiber.Ctx) error { return nil }

type gigPublisherStub struct{}

func (gigPublisherStub) CreateDraft(context.Context, gateway.CreateGigDraftRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigPublisherStub) UpdateBasicInfo(context.Context, gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigPublisherStub) ReplacePackages(context.Context, gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigPublisherStub) ReplaceQuestions(context.Context, gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigPublisherStub) ReplaceMedia(context.Context, gateway.ReplaceGigMediaRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigPublisherStub) GetDraft(context.Context, gateway.GetGigDraftRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigPublisherStub) GetBySlug(context.Context, gateway.GetGigBySlugRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigPublisherStub) Publish(context.Context, gateway.PublishGigRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

type reviewClientStub struct{}

func (reviewClientStub) CreateReview(context.Context, gateway.CreateReviewRequest) (*gateway.CreateReviewResult, error) {
	return &gateway.CreateReviewResult{ReviewID: "review-1"}, nil
}

func (reviewClientStub) GetGigReviews(context.Context, gateway.GetGigReviewsRequest) (*gateway.GetGigReviewsResult, error) {
	return &gateway.GetGigReviewsResult{}, nil
}

func (reviewClientStub) GetGigReviewsSummary(context.Context, gateway.GetGigReviewsSummaryRequest) (*gateway.ReviewSummary, error) {
	return &gateway.ReviewSummary{}, nil
}

func (reviewClientStub) GetUserRatingSummaryByUsername(context.Context, gateway.GetUserRatingSummaryByUsernameRequest) (*gateway.ReviewSummary, error) {
	return &gateway.ReviewSummary{}, nil
}

func (reviewClientStub) Close() error { return nil }

type userClientStub struct{}

func (userClientStub) GetDetailedUserByUsername(context.Context, string) (*gateway.User, error) {
	return &gateway.User{UserID: "user-1", Username: "alex"}, nil
}

func (userClientStub) Close() error { return nil }

type gigServiceStub struct{}

func (gigServiceStub) CreateDraft(context.Context, gateway.CreateGigDraftRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigServiceStub) UpdateBasicInfo(context.Context, gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigServiceStub) ReplacePackages(context.Context, gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigServiceStub) ReplaceQuestions(context.Context, gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigServiceStub) ReplaceMedia(context.Context, gateway.ReplaceGigMediaRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigServiceStub) GetDraft(context.Context, gateway.GetGigDraftRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigServiceStub) GetBySlug(context.Context, gateway.GetGigBySlugRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
}

func (gigServiceStub) Publish(context.Context, gateway.PublishGigRequest) (*gateway.Gig, error) {
	return &gateway.Gig{GigID: "gig-1"}, nil
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
			JWT:              config.JWTConfig{AccessSecret: "local-dev-access-secret-change-me"},
			RegistrationSaga: config.RegistrationSagaConfig{Address: "127.0.0.1:9500"},
		}
	})

	It("loads config from the current environment", func() {
		prevWD, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			Expect(os.Chdir(prevWD)).To(Succeed())
			Expect(os.Unsetenv("HTTP_PORT")).To(Succeed())
			Expect(os.Unsetenv("PAYMENT_SERVICE_ADDRESS")).To(Succeed())
		}()

		tmpDir := GinkgoT().TempDir()
		Expect(os.Chdir(tmpDir)).To(Succeed())
		Expect(os.Setenv("HTTP_PORT", "9099")).To(Succeed())
		Expect(os.Setenv("PAYMENT_SERVICE_ADDRESS", "127.0.0.1:9506")).To(Succeed())

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

	It("constructs and syncs the logger", func() {
		lc := &lifecycleStub{}
		logger, err := ProvideLogger(lc, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(logger).NotTo(BeNil())
		Expect(lc.hooks).To(HaveLen(1))
		Expect(func() { _ = lc.hooks[0].OnStop(context.Background()) }).NotTo(Panic())
	})

	It("constructs the registration application service", func() {
		svc, err := ProvideRegistrationService(registrationPublisherStub{}, tokenIssuerStub{}, lg)

		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())
	})

	It("constructs the HTTP server", func() {
		httpSrv, err := ProvideHTTPServer(cfg, lg)
		Expect(err).NotTo(HaveOccurred())
		Expect(httpSrv).NotTo(BeNil())
	})

	It("registers the HTTP server lifecycle hook", func() {
		lc := &lifecycleStub{}
		srv := &httpServerStub{app: fiber.New()}

		InvokeRunHTTPServer(lc, srv, lg)

		Expect(lc.hooks).To(HaveLen(1))
		Expect(lc.hooks[0].OnStart(context.Background())).To(Succeed())
		Expect(func() { _ = lc.hooks[0].OnStop(context.Background()) }).NotTo(Panic())
		Eventually(func() int { return srv.startCalls }).Should(Equal(1))
		Eventually(func() int { return srv.shutdownCalls }).Should(Equal(1))
	})

	It("constructs the v1 auth handler", func() {
		handler, err := ProvideHTTPV1AuthHandler(registrationServiceStub{}, authSessionServiceStub{}, lg)

		Expect(err).NotTo(HaveOccurred())
		Expect(handler).NotTo(BeNil())
	})

	It("constructs the v1 gig handler", func() {
		handler, err := ProvideHTTPV1GigHandler(cfg, gigServiceStub{}, lg)

		Expect(err).NotTo(HaveOccurred())
		Expect(handler).NotTo(BeNil())
	})

	It("registers versioned routes", func() {
		srv := &httpServerStub{app: fiber.New()}
		handler := &authHandlerStub{}
		gigHandler := &gigHandlerStub{}
		orderHandler := &orderHandlerStub{}
		reviewHandler := &reviewHandlerStub{}
		searchHandler := &searchHandlerStub{}
		onboardingHandler := &onboardingHandlerStub{}
		InvokeRegisterHTTPV1Routes(srv, handler, gigHandler, orderHandler, reviewHandler, searchHandler, onboardingHandler)

		Expect(handler.registered).To(BeTrue())
		Expect(gigHandler.registered).To(BeTrue())
		Expect(orderHandler.registered).To(BeTrue())
		Expect(reviewHandler.registered).To(BeTrue())
		Expect(onboardingHandler.registered).To(BeTrue())
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

	It("constructs a token issuer and registers its shutdown hook", func() {
		server, address := startAuthServer()
		defer server.Stop()

		lc := &lifecycleStub{}
		cfg.AuthService.Address = address

		tok, err := ProvideTokenIssuer(lc, cfg, lg)
		Expect(err).NotTo(HaveOccurred())
		Expect(tok).NotTo(BeNil())
		Expect(lc.hooks).To(HaveLen(1))
		Expect(lc.hooks[0].OnStop(context.Background())).To(Succeed())
	})

	It("constructs the gig publisher and gig application service", func() {
		lc := &lifecycleStub{}
		cfg.GigService.Address = "127.0.0.1:1"

		pub, err := ProvideGigPublisher(lc, cfg, lg)
		Expect(err).NotTo(HaveOccurred())
		Expect(pub).NotTo(BeNil())
		Expect(lc.hooks).To(HaveLen(1))
		Expect(lc.hooks[0].OnStop(context.Background())).To(Succeed())

		svc, err := ProvideGigService(gigPublisherStub{}, reviewClientStub{}, userClientStub{}, lg)
		Expect(err).NotTo(HaveOccurred())
		Expect(svc).NotTo(BeNil())
	})

	It("returns grpc dial errors from registration publisher construction", func() {
		lc := &lifecycleStub{}
		cfg.RegistrationSaga.Address = ""

		pub, err := ProvideRegistrationPublisher(lc, cfg, lg)

		Expect(pub).To(BeNil())
		Expect(err).To(MatchError(grpcclient.ErrEmptyAddress))
		Expect(lc.hooks).To(BeEmpty())
	})

	It("returns grpc dial errors from token issuer construction", func() {
		lc := &lifecycleStub{}
		cfg.AuthService.Address = ""

		tok, err := ProvideTokenIssuer(lc, cfg, lg)
		Expect(tok).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(lc.hooks).To(BeEmpty())
	})

	It("returns grpc dial errors from gig publisher construction", func() {
		lc := &lifecycleStub{}
		cfg.GigService.Address = ""

		pub, err := ProvideGigPublisher(lc, cfg, lg)
		Expect(pub).To(BeNil())
		Expect(err).To(MatchError(grpcclient.ErrEmptyGigServiceAddress))
		Expect(lc.hooks).To(BeEmpty())
	})

	It("emits the start log without panicking", func() {
		Expect(func() { InvokeStartLog(lg) }).NotTo(Panic())
	})
})

func startRegistrationServer() (*grpc.Server, string) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		Skip("gRPC sockets are unavailable in this environment")
		return nil, ""
	}

	server := grpc.NewServer()
	registrationv1.RegisterRegistrationServiceServer(server, fakeRegistrationServer{})

	go func() {
		_ = server.Serve(lis)
	}()

	return server, lis.Addr().String()
}

func startAuthServer() (*grpc.Server, string) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		Skip("gRPC sockets are unavailable in this environment")
		return nil, ""
	}

	server := grpc.NewServer()
	authv1.RegisterAuthQueryServiceServer(server, authServerStub{})

	go func() {
		_ = server.Serve(lis)
	}()

	return server, lis.Addr().String()
}

type authServerStub struct {
	authv1.UnimplementedAuthQueryServiceServer
}

func (authServerStub) IssueRegistrationTokens(context.Context, *authv1.IssueRegistrationTokensRequest) (*authv1.IssueRegistrationTokensResponse, error) {
	return &authv1.IssueRegistrationTokensResponse{TokenType: "Bearer"}, nil
}
