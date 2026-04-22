package http

import (
	"api-gateway/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	gateway "api-gateway/internal/domain"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestHTTP(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP Suite")
}

var _ = Describe("AuthHandler", func() {
	var (
		ctrl    *gomock.Controller
		service *MockRegistrationService
		logger  logging.Logger
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		service = NewMockRegistrationService(ctrl)

		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("NewAuthHandler", func() {
		It("validates nil collaborators", func() {
			handler, err := NewAuthHandler(nil, logger)
			Expect(handler).To(BeNil())
			Expect(err).To(MatchError(ErrNilRegistrationService))

			handler, err = NewAuthHandler(service, nil)
			Expect(handler).To(BeNil())
			Expect(err).To(MatchError(ErrNilLogger))
		})
	})

	Describe("signup endpoint", func() {
		var (
			app    *fiber.App
			client *resty.Client
			stop   func()
		)

		BeforeEach(func() {
			handler, err := NewAuthHandler(service, logger)
			Expect(err).NotTo(HaveOccurred())

			app = fiber.New()
			v1 := app.Group("/v1")
			handler.RegisterRoutes(v1)

			client, stop = startFiberClient(app)
		})

		AfterEach(func() {
			stop()
		})

		It("rejects an invalid request body", func() {
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody("{").
				Post("/v1/auth/sign-up")

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(fiber.StatusBadRequest))
			Expect(resp.String()).To(MatchJSON(`{"error":"invalid request body"}`))
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

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]any{
					"email":     "alex@example.com",
					"username":  "alex",
					"password":  "password123",
					"firstName": "Alex",
					"surname":   "Doe",
				}).
				Post("/v1/auth/sign-up")

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(fiber.StatusAccepted))
			Expect(resp.String()).To(MatchJSON(`{"session_id":"session-1","status":"pending"}`))
		})

		It("maps invalid email to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidEmail)

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]any{
					"email":    "alex@example.com",
					"username": "alex",
					"password": "password123",
				}).
				Post("/v1/auth/sign-up")

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(fiber.StatusBadRequest))
			Expect(resp.String()).To(MatchJSON(`{"error":"invalid email"}`))
		})

		It("maps invalid username to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidUsername)

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]any{
					"email":    "alex@example.com",
					"username": "alex",
					"password": "password123",
				}).
				Post("/v1/auth/sign-up")

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(fiber.StatusBadRequest))
			Expect(resp.String()).To(MatchJSON(`{"error":"invalid username"}`))
		})

		It("maps invalid password to 400", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrInvalidPassword)

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]any{
					"email":    "alex@example.com",
					"username": "alex",
					"password": "password123",
				}).
				Post("/v1/auth/sign-up")

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(fiber.StatusBadRequest))
			Expect(resp.String()).To(MatchJSON(`{"error":"invalid password"}`))
		})

		It("maps registration conflicts to 409", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, &gateway.RegistrationConflictError{
					State:         "found_completed",
					UsernameTaken: true,
					EmailTaken:    false,
				})

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]any{
					"email":    "alex@example.com",
					"username": "alex",
					"password": "password123",
				}).
				Post("/v1/auth/sign-up")

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(fiber.StatusConflict))

			var body map[string]any
			Expect(json.Unmarshal(resp.Body(), &body)).To(Succeed())
			Expect(body["error"]).To(Equal("registration conflict"))
			Expect(body["state"]).To(Equal("found_completed"))
			Expect(body["username_taken"]).To(Equal(true))
			Expect(body["email_taken"]).To(Equal(false))
		})

		It("maps internal failures to 500", func() {
			service.EXPECT().
				SignUp(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("boom"))

			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(map[string]any{
					"email":    "alex@example.com",
					"username": "alex",
					"password": "password123",
				}).
				Post("/v1/auth/sign-up")

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(fiber.StatusInternalServerError))
			Expect(resp.String()).To(MatchJSON(`{"error":"internal server error"}`))
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

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())
		defer ln.Close()

		done := make(chan error, 1)
		go func() {
			done <- srv.app.Listener(ln)
		}()

		Expect(srv.Shutdown(context.Background())).To(Succeed())
	})

	It("starts serving over the configured address", func() {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())
		addr := ln.Addr().(*net.TCPAddr)
		Expect(ln.Close()).To(Succeed())

		srvAny, err := NewServer(serverConfig(addr.Port), logger)
		Expect(err).NotTo(HaveOccurred())
		srv := srvAny.(*server)

		done := make(chan error, 1)
		go func() {
			done <- srv.Start()
		}()

		client := resty.New().
			SetBaseURL("http://127.0.0.1:" + fmt.Sprint(addr.Port)).
			SetRetryCount(10).
			SetRetryWaitTime(20 * time.Millisecond)

		Eventually(func() int {
			resp, reqErr := client.R().Get("/not-found")
			if reqErr != nil {
				return 0
			}
			return resp.StatusCode()
		}, time.Second, 20*time.Millisecond).Should(Equal(fiber.StatusNotFound))

		Expect(srv.Shutdown(context.Background())).To(Succeed())
	})
})

var _ = Describe("swagger doc anchor", func() {
	It("is callable for coverage", func() {
		Expect(func() { swaggerSignUpDoc() }).NotTo(Panic())
	})
})

func startFiberClient(app *fiber.App) (*resty.Client, func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	Expect(err).NotTo(HaveOccurred())

	done := make(chan error, 1)
	go func() {
		done <- app.Listener(ln)
	}()

	client := resty.New().
		SetBaseURL("http://" + ln.Addr().String()).
		SetRetryCount(3).
		SetRetryWaitTime(20 * time.Millisecond)

	return client, func() {
		Expect(app.Shutdown()).To(Succeed())
		Eventually(done, time.Second).Should(Receive(BeNil()))
	}
}

func serverConfig(port int) config.HTTPConfig {
	return config.HTTPConfig{
		Host: "127.0.0.1",
		Port: port,
	}
}
