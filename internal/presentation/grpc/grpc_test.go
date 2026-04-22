package grpc

import (
	"context"
	"net"
	"testing"

	gateway "api-gateway/internal/domain"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	registrationv1 "github.com/ofm-microseervices/ofm-common/proto/registration/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPC(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Suite")
}

type fakeRegistrationServer struct {
	registrationv1.UnimplementedRegistrationServiceServer
	startFn func(context.Context, *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error)
}

func (f *fakeRegistrationServer) StartRegistration(ctx context.Context, req *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error) {
	return f.startFn(ctx, req)
}

var _ = Describe("registration mapper", func() {
	It("maps the gateway signup request to the proto request", func() {
		mapr := newRegistrationMapper()

		req := mapr.ToStartRegistrationRequest(gateway.SignUpRequest{
			ClientID:  "client-1",
			Email:     "alex@example.com",
			Username:  "alex",
			Password:  "password123",
			FirstName: "Alex",
			Surname:   "Doe",
		})

		Expect(req.ClientId).To(Equal("client-1"))
		Expect(req.Email).To(Equal("alex@example.com"))
		Expect(req.Username).To(Equal("alex"))
		Expect(req.Password).To(Equal("password123"))
		Expect(req.FirstName).To(Equal("Alex"))
		Expect(req.Surname).To(Equal("Doe"))
	})

	It("maps proto responses back to the gateway result", func() {
		mapr := newRegistrationMapper()

		result := mapr.ToSignUpResult(&registrationv1.StartRegistrationResponse{
			SessionId:     "session-1",
			ClientId:      "client-1",
			UserId:        "user-1",
			Status:        "conflict",
			ConflictState: "found_completed",
			UsernameTaken: true,
			EmailTaken:    false,
		})

		Expect(result).To(Equal(&gateway.SignUpResult{
			SessionID:     "session-1",
			ClientID:      "client-1",
			UserID:        "user-1",
			Status:        "conflict",
			ConflictState: "found_completed",
			UsernameTaken: true,
			EmailTaken:    false,
		}))
	})

	It("maps grpc invalid-argument errors to gateway validation errors", func() {
		mapr := newRegistrationMapper()

		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "invalid email"))).To(MatchError(gateway.ErrInvalidEmail))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "invalid username"))).To(MatchError(gateway.ErrInvalidUsername))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "invalid password"))).To(MatchError(gateway.ErrInvalidPassword))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "something else"))).To(MatchError(gateway.ErrFailedToStartRegistration))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.Internal, "boom"))).To(MatchError(gateway.ErrFailedToStartRegistration))
		Expect(mapr.ToStartRegistrationError(context.Canceled)).To(MatchError(gateway.ErrFailedToStartRegistration))
	})
})

var _ = Describe("registration client", func() {
	var logger logging.Logger

	BeforeEach(func() {
		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	It("validates constructor input", func() {
		cl, err := NewClient(RegistrationSagaConfig{}, logger)
		Expect(cl).To(BeNil())
		Expect(err).To(MatchError(ErrEmptyAddress))

		cl, err = NewClient(RegistrationSagaConfig{Address: "127.0.0.1:9090"}, nil)
		Expect(cl).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("returns success results from a real grpc server", func() {
		server, address := startRegistrationServer(func(_ context.Context, req *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error) {
			Expect(req.Email).To(Equal("alex@example.com"))
			Expect(req.FirstName).To(Equal("Alex"))
			return &registrationv1.StartRegistrationResponse{
				SessionId: "session-1",
				Status:    "pending",
			}, nil
		})
		defer server.Stop()

		cl, err := NewClient(RegistrationSagaConfig{Address: address}, logger)
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			Expect(cl.Close()).To(Succeed())
		}()

		result, err := cl.StartRegistration(context.Background(), SignUpRequest{
			Email:     "alex@example.com",
			Username:  "alex",
			Password:  "password123",
			FirstName: "Alex",
			Surname:   "Doe",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&gateway.SignUpResult{
			SessionID: "session-1",
			Status:    "pending",
		}))
	})

	It("maps validation failures returned by grpc", func() {
		server, address := startRegistrationServer(func(context.Context, *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error) {
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		})
		defer server.Stop()

		cl, err := NewClient(RegistrationSagaConfig{Address: address}, logger)
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			Expect(cl.Close()).To(Succeed())
		}()

		result, err := cl.StartRegistration(context.Background(), SignUpRequest{
			Email:    "alex@example.com",
			Username: "alex",
			Password: "password123",
		})

		Expect(result).To(BeNil())
		Expect(err).To(MatchError(gateway.ErrInvalidEmail))
	})

	It("closes a nil client safely", func() {
		var cl *client
		Expect(cl.Close()).To(Succeed())
	})

	It("closes a client with a nil connection safely", func() {
		cl := &client{log: logger}
		Expect(cl.Close()).To(Succeed())
	})
})

func startRegistrationServer(
	startFn func(context.Context, *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error),
) (*grpc.Server, string) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	Expect(err).NotTo(HaveOccurred())

	server := grpc.NewServer()
	registrationv1.RegisterRegistrationServiceServer(server, &fakeRegistrationServer{startFn: startFn})

	go func() {
		_ = server.Serve(lis)
	}()

	return server, lis.Addr().String()
}
