package http

import (
	"api-gateway/config"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gateway "api-gateway/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestHTTP(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Suite")
}

type authSessionServiceStub struct{}

func (authSessionServiceStub) SignIn(context.Context, gateway.SignInRequest) (*gateway.AuthTokensResult, error) {
	return &gateway.AuthTokensResult{}, nil
}

func (authSessionServiceStub) Refresh(context.Context, gateway.RefreshTokensRequest) (*gateway.AuthTokensResult, error) {
	return &gateway.AuthTokensResult{}, nil
}

func (authSessionServiceStub) Close() error { return nil }

var _ = Describe("AuthHandler", func() {
	var (
		ctrl    *gomock.Controller
		service *MockRegistrationService
		session authSessionServiceStub
		logger  logging.Logger
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		service = NewMockRegistrationService(ctrl)
		session = authSessionServiceStub{}

		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("NewAuthHandler", func() {
		It("validates nil collaborators", func() {
			handler, err := NewAuthHandler(nil, session, logger)
			Expect(handler).To(BeNil())
			Expect(err).To(MatchError(ErrNilRegistrationService))

			handler, err = NewAuthHandler(service, nil, logger)
			Expect(handler).To(BeNil())
			Expect(err).To(MatchError(ErrNilAuthSessionService))

			handler, err = NewAuthHandler(service, session, nil)
			Expect(handler).To(BeNil())
			Expect(err).To(MatchError(ErrNilLogger))
		})
	})

	Describe("signup endpoint", func() {
		var (
			app *fiber.App
		)

		BeforeEach(func() {
			handler, err := NewAuthHandler(service, session, logger)
			Expect(err).NotTo(HaveOccurred())

			app = fiber.New()
			v1 := app.Group("/v1")
			handler.RegisterRoutes(v1)
		})

		It("rejects an invalid request body", func() {
			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", "{"), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"` + ErrInvalidRequestBody.Error() + `"}`))
		})

		It("returns 202 for an accepted signup request", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gateway.SignUpRequest{
					Email:     "alex@example.com",
					Username:  "alex",
					Password:  "password123",
					FirstName: "Alex",
					Surname:   "Doe",
				}).
				Return(&gateway.SignUpResult{
					SessionID: "session-1",
					Status:    "pending",
				}, nil)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123","firstName":"Alex","surname":"Doe"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusAccepted))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"session_id":"session-1","status":"pending"}`))
		})

		It("maps invalid email to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidEmail)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"invalid email"}`))
		})

		It("maps invalid session id to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidSessionID)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"invalid session id"}`))
		})

		It("maps invalid client id to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidClientID)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"invalid client id"}`))
		})

		It("maps invalid verification code to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidVerificationCode)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"invalid verification code"}`))
		})

		It("maps invalid username to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidUsername)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"invalid username"}`))
		})

		It("maps invalid password to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidPassword)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"invalid password"}`))
		})

		It("maps registration conflicts to 409", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, &gateway.RegistrationConflictError{
					State:         "found_completed",
					UsernameTaken: true,
					EmailTaken:    false,
				})

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusConflict))

			var body map[string]any
			raw, _ := io.ReadAll(resp.Body)
			Expect(json.Unmarshal(raw, &body)).To(Succeed())
			Expect(body["error"]).To(Equal("registration conflict"))
			Expect(body["state"]).To(Equal("found_completed"))
			Expect(body["username_taken"]).To(Equal(true))
			Expect(body["email_taken"]).To(Equal(false))
		})

		It("maps internal failures to 500", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom"))

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up", `{"email":"alex@example.com","username":"alex","password":"password123"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusInternalServerError))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"internal server error"}`))
		})
	})

	Describe("verify-email endpoint", func() {
		var app *fiber.App

		BeforeEach(func() {
			handler, err := NewAuthHandler(service, session, logger)
			Expect(err).NotTo(HaveOccurred())
			app = fiber.New()
			v1 := app.Group("/v1")
			handler.RegisterRoutes(v1)
		})

		It("returns 202 for accepted verification", func() {
			service.EXPECT().
				VerifyEmail(gomock.Any(), gateway.VerifyEmailRequest{
					SessionID: "session-1",
					ClientID:  "client-1",
					Code:      "123456",
				}).
				Return(&gateway.VerifyEmailResult{SessionID: "session-1", ClientID: "client-1", Status: "verifying_email"}, nil)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/verify-email", `{"session_id":"session-1","client_id":"client-1","code":"123456"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusAccepted))
		})

		It("rejects an invalid request body", func() {
			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/verify-email", "{"), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"` + ErrInvalidRequestBody.Error() + `"}`))
		})

		It("maps invalid client id to 400", func() {
			service.EXPECT().VerifyEmail(gomock.Any(), gomock.Any()).Return(nil, gateway.ErrInvalidClientID)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/verify-email", `{"session_id":"session-1","client_id":"client-1","code":"123456"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		})

		It("maps invalid session id to 400", func() {
			service.EXPECT().VerifyEmail(gomock.Any(), gomock.Any()).Return(nil, gateway.ErrInvalidSessionID)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/verify-email", `{"session_id":"session-1","client_id":"client-1","code":"123456"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		})

		It("maps invalid verification code to 400", func() {
			service.EXPECT().VerifyEmail(gomock.Any(), gomock.Any()).Return(nil, gateway.ErrInvalidVerificationCode)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/verify-email", `{"session_id":"session-1","client_id":"client-1","code":"123456"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		})

		It("maps internal failures to 500", func() {
			service.EXPECT().VerifyEmail(gomock.Any(), gomock.Any()).Return(nil, errors.New("boom"))

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/verify-email", `{"session_id":"session-1","client_id":"client-1","code":"123456"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusInternalServerError))
		})
	})

	Describe("complete-registration endpoint", func() {
		var app *fiber.App

		BeforeEach(func() {
			handler, err := NewAuthHandler(service, session, logger)
			Expect(err).NotTo(HaveOccurred())
			app = fiber.New()
			v1 := app.Group("/v1")
			handler.RegisterRoutes(v1)
		})

		It("returns tokens for completed registration", func() {
			service.EXPECT().
				CompleteRegistration(gomock.Any(), gateway.CompleteRegistrationRequest{
					SessionID: "session-1",
					ClientID:  "client-1",
				}).
				Return(&gateway.CompleteRegistrationResult{
					UserID:      "user-1",
					AccessToken: "access",
					TokenType:   "Bearer",
				}, nil)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/complete", `{"session_id":"session-1","client_id":"client-1"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})

		It("rejects an invalid request body", func() {
			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/complete", "{"), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
			body, _ := io.ReadAll(resp.Body)
			Expect(string(body)).To(MatchJSON(`{"error":"` + ErrInvalidRequestBody.Error() + `"}`))
		})

		It("maps already-claimed errors to 410", func() {
			service.EXPECT().CompleteRegistration(gomock.Any(), gomock.Any()).Return(nil, gateway.ErrRegistrationAlreadyClaimed)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/complete", `{"session_id":"session-1","client_id":"client-1"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusGone))
		})

		It("maps not-completed errors to 409", func() {
			service.EXPECT().CompleteRegistration(gomock.Any(), gomock.Any()).Return(nil, gateway.ErrRegistrationNotCompleted)

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/complete", `{"session_id":"session-1","client_id":"client-1"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusConflict))
		})

		It("maps internal failures to 500", func() {
			service.EXPECT().CompleteRegistration(gomock.Any(), gomock.Any()).Return(nil, errors.New("boom"))

			resp, err := app.Test(jsonRequest("POST", "/v1/auth/sign-up/complete", `{"session_id":"session-1","client_id":"client-1"}`), -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusInternalServerError))
		})
	})
})

var _ = Describe("HTTP server", func() {
	var logger logging.Logger

	BeforeEach(func() {
		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	It("validates a nil logger", func() {
		srv, err := NewServer(serverConfig(0), nil)
		Expect(srv).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("exposes the underlying fiber app and shuts down cleanly", func() {
		srvAny, err := NewServer(serverConfig(0), logger)
		Expect(err).NotTo(HaveOccurred())

		srv := srvAny.(*server)
		Expect(srv.App()).NotTo(BeNil())

		Expect(srv.Shutdown(context.Background())).To(Succeed())
	})

	It("starts the server on a configured address", func() {
		srvAny, err := NewServer(serverConfig(0), logger)
		Expect(err).NotTo(HaveOccurred())

		srv := srvAny.(*server)
		errCh := make(chan error, 1)
		go func() { errCh <- srv.Start() }()
		time.Sleep(100 * time.Millisecond)
		Expect(srv.Shutdown(context.Background())).To(Succeed())
		Eventually(errCh, 3*time.Second).Should(Receive())
	})
})

var _ = Describe("swagger doc anchor", func() {
	It("is callable for coverage", func() {
		Expect(func() { swaggerSignUpDoc() }).NotTo(Panic())
	})
})

func serverConfig(port int) config.HTTPConfig {
	return config.HTTPConfig{
		Host: "127.0.0.1",
		Port: port,
	}
}

func jsonRequest(method, target, body string) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
