package service

import (
	"context"
	"errors"
	"testing"

	gateway "api-gateway/internal/domain"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestApplication(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Suite")
}

var _ = Describe("RegistrationService", func() {
	var (
		ctrl      *gomock.Controller
		publisher *MockRegistrationPublisher
		logger    logging.Logger
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		publisher = NewMockRegistrationPublisher(ctrl)

		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("New", func() {
		It("validates nil collaborators", func() {
			svc, err := New(nil, logger)
			Expect(svc).To(BeNil())
			Expect(err).To(MatchError(ErrNilRegistrationClient))

			svc, err = New(publisher, nil)
			Expect(svc).To(BeNil())
			Expect(err).To(MatchError(ErrNilLogger))
		})
	})

	Describe("SignUp", func() {
		It("rejects an invalid email", func() {
			svc, err := New(publisher, logger)
			Expect(err).NotTo(HaveOccurred())

			result, err := svc.SignUp(context.Background(), gateway.SignUpRequest{
				Email:    "bad-email",
				Username: "alex",
				Password: "password123",
			})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError(gateway.ErrInvalidEmail))
		})

		It("rejects an empty username", func() {
			svc, err := New(publisher, logger)
			Expect(err).NotTo(HaveOccurred())

			result, err := svc.SignUp(context.Background(), gateway.SignUpRequest{
				Email:    "alex@example.com",
				Username: "   ",
				Password: "password123",
			})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError(gateway.ErrInvalidUsername))
		})

		It("rejects a short password", func() {
			svc, err := New(publisher, logger)
			Expect(err).NotTo(HaveOccurred())

			result, err := svc.SignUp(context.Background(), gateway.SignUpRequest{
				Email:    "alex@example.com",
				Username: "alex",
				Password: "short",
			})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError(gateway.ErrInvalidPassword))
		})

		It("trims the public payload before delegating to the publisher", func() {
			svc, err := New(publisher, logger)
			Expect(err).NotTo(HaveOccurred())

			publisher.EXPECT().
				StartRegistration(gomock.Any(), gateway.SignUpRequest{
					ClientID:  "client-1",
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

			result, err := svc.SignUp(context.Background(), gateway.SignUpRequest{
				ClientID:  " client-1 ",
				Email:     "  alex@example.com ",
				Username:  " alex ",
				Password:  "password123",
				FirstName: " Alex ",
				Surname:   " Doe ",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(&gateway.SignUpResult{
				SessionID: "session-1",
				Status:    "pending",
			}))
		})

		It("returns publisher failures unchanged", func() {
			svc, err := New(publisher, logger)
			Expect(err).NotTo(HaveOccurred())

			publisher.EXPECT().
				StartRegistration(gomock.Any(), gomock.Any()).
				Return(nil, gateway.ErrFailedToStartRegistration)

			result, err := svc.SignUp(context.Background(), gateway.SignUpRequest{
				Email:    "alex@example.com",
				Username: "alex",
				Password: "password123",
			})

			Expect(result).To(BeNil())
			Expect(err).To(MatchError(gateway.ErrFailedToStartRegistration))
		})

		It("maps conflict responses to a conflict error", func() {
			svc, err := New(publisher, logger)
			Expect(err).NotTo(HaveOccurred())

			publisher.EXPECT().
				StartRegistration(gomock.Any(), gomock.Any()).
				Return(&gateway.SignUpResult{
					Status:        "conflict",
					ConflictState: "found_completed",
					UsernameTaken: true,
					EmailTaken:    false,
				}, nil)

			result, err := svc.SignUp(context.Background(), gateway.SignUpRequest{
				Email:    "alex@example.com",
				Username: "alex",
				Password: "password123",
			})

			Expect(result).To(BeNil())
			var conflictErr *gateway.RegistrationConflictError
			Expect(errors.As(err, &conflictErr)).To(BeTrue())
			Expect(conflictErr.State).To(Equal("found_completed"))
			Expect(conflictErr.UsernameTaken).To(BeTrue())
			Expect(conflictErr.EmailTaken).To(BeFalse())
		})
	})
})
