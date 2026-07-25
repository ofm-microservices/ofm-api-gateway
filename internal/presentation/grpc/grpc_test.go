package grpc

import (
	"context"
	"net"
	"testing"

	gateway "api-gateway/internal/domain"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	authv1 "github.com/ofm-microservices/ofm-common/proto/auth/v1"
	registrationv1 "github.com/ofm-microservices/ofm-common/proto/registration/v1"
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
	startFn  func(context.Context, *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error)
	verifyFn func(context.Context, *registrationv1.VerifyEmailRequest) (*registrationv1.VerifyEmailResponse, error)
	statusFn func(context.Context, *registrationv1.GetRegistrationStatusRequest) (*registrationv1.GetRegistrationStatusResponse, error)
}

func (f *fakeRegistrationServer) StartRegistration(ctx context.Context, req *registrationv1.StartRegistrationRequest) (*registrationv1.StartRegistrationResponse, error) {
	if f.startFn == nil {
		return &registrationv1.StartRegistrationResponse{Status: "pending"}, nil
	}
	return f.startFn(ctx, req)
}

func (f *fakeRegistrationServer) VerifyEmail(ctx context.Context, req *registrationv1.VerifyEmailRequest) (*registrationv1.VerifyEmailResponse, error) {
	if f.verifyFn == nil {
		return &registrationv1.VerifyEmailResponse{Status: "verifying_email"}, nil
	}
	return f.verifyFn(ctx, req)
}

func (f *fakeRegistrationServer) GetRegistrationStatus(ctx context.Context, req *registrationv1.GetRegistrationStatusRequest) (*registrationv1.GetRegistrationStatusResponse, error) {
	if f.statusFn == nil {
		return &registrationv1.GetRegistrationStatusResponse{Status: "completed", UserId: "user-1"}, nil
	}
	return f.statusFn(ctx, req)
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

	It("maps verify-email requests and results", func() {
		mapr := newRegistrationMapper()

		req := mapr.ToVerifyEmailRequest(gateway.VerifyEmailRequest{
			SessionID: "session-1",
			ClientID:  "client-1",
			Code:      "123456",
		})
		Expect(req.GetSessionId()).To(Equal("session-1"))
		Expect(req.GetClientId()).To(Equal("client-1"))
		Expect(req.GetCode()).To(Equal("123456"))

		result := mapr.ToVerifyEmailResult(&registrationv1.VerifyEmailResponse{
			SessionId: "session-1",
			ClientId:  "client-1",
			Status:    "verifying_email",
		})
		Expect(result).To(Equal(&gateway.VerifyEmailResult{SessionID: "session-1", ClientID: "client-1", Status: "verifying_email"}))
	})

	It("maps registration status and token results", func() {
		mapr := newRegistrationMapper()

		status := mapr.ToRegistrationStatus(&registrationv1.GetRegistrationStatusResponse{
			SessionId: "session-1",
			ClientId:  "client-1",
			UserId:    "user-1",
			Status:    "completed",
		})
		Expect(status).To(Equal(&gateway.RegistrationStatus{SessionID: "session-1", ClientID: "client-1", UserID: "user-1", Status: "completed"}))

		tokens := mapr.ToCompleteRegistrationResult(&authv1.IssueRegistrationTokensResponse{
			UserId:       "user-1",
			AccessToken:  "access",
			RefreshToken: "refresh",
			TokenType:    "Bearer",
			ExpiresIn:    900,
		})
		Expect(tokens).To(Equal(&gateway.CompleteRegistrationResult{
			UserID:       "user-1",
			AccessToken:  "access",
			RefreshToken: "refresh",
			TokenType:    "Bearer",
			ExpiresIn:    900,
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

	It("maps the remaining registration start error branches", func() {
		mapr := newRegistrationMapper()
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "invalid session id"))).To(MatchError(gateway.ErrInvalidSessionID))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "invalid client id"))).To(MatchError(gateway.ErrInvalidClientID))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "invalid verification code"))).To(MatchError(gateway.ErrInvalidVerificationCode))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "registration is not ready for this operation"))).To(MatchError(gateway.ErrRegistrationNotCompleted))
		Expect(mapr.ToStartRegistrationError(status.Error(codes.InvalidArgument, "registration session not found"))).To(MatchError(gateway.ErrRegistrationNotCompleted))
	})

	It("maps registration-status errors to completion errors", func() {
		mapr := newRegistrationMapper()

		Expect(mapr.ToRegistrationStatusError(status.Error(codes.InvalidArgument, "invalid session id"))).To(MatchError(gateway.ErrInvalidSessionID))
		Expect(mapr.ToRegistrationStatusError(status.Error(codes.InvalidArgument, "invalid client id"))).To(MatchError(gateway.ErrInvalidClientID))
		Expect(mapr.ToRegistrationStatusError(status.Error(codes.InvalidArgument, "registration is not ready for this operation"))).To(MatchError(gateway.ErrRegistrationNotCompleted))
		Expect(mapr.ToRegistrationStatusError(status.Error(codes.InvalidArgument, "registration session not found"))).To(MatchError(gateway.ErrRegistrationNotCompleted))
		Expect(mapr.ToRegistrationStatusError(status.Error(codes.Internal, "boom"))).To(MatchError(gateway.ErrFailedToCompleteRegistration))
		Expect(mapr.ToRegistrationStatusError(context.Canceled)).To(MatchError(gateway.ErrFailedToCompleteRegistration))
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

		cl, err = NewClient(RegistrationSagaConfig{Address: "127.0.0.1:9500"}, nil)
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

	It("maps verification and status grpc failures", func() {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			Skip("gRPC sockets are unavailable in this environment")
		}

		server := grpc.NewServer()
		registrationv1.RegisterRegistrationServiceServer(server, &fakeRegistrationServer{
			verifyFn: func(context.Context, *registrationv1.VerifyEmailRequest) (*registrationv1.VerifyEmailResponse, error) {
				return nil, status.Error(codes.InvalidArgument, "invalid email")
			},
			statusFn: func(context.Context, *registrationv1.GetRegistrationStatusRequest) (*registrationv1.GetRegistrationStatusResponse, error) {
				return nil, status.Error(codes.Internal, "boom")
			},
		})
		go func() { _ = server.Serve(lis) }()
		defer server.Stop()

		cl, err := NewClient(RegistrationSagaConfig{Address: lis.Addr().String()}, logger)
		Expect(err).NotTo(HaveOccurred())
		defer func() { Expect(cl.Close()).To(Succeed()) }()

		_, err = cl.VerifyEmail(context.Background(), VerifyEmailRequest{
			SessionID: "session-1",
			ClientID:  "client-1",
			Code:      "123456",
		})
		Expect(err).To(MatchError(gateway.ErrInvalidEmail))
	})

	It("maps registration-status grpc failures", func() {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			Skip("gRPC sockets are unavailable in this environment")
		}

		server := grpc.NewServer()
		registrationv1.RegisterRegistrationServiceServer(server, &fakeRegistrationServer{
			statusFn: func(context.Context, *registrationv1.GetRegistrationStatusRequest) (*registrationv1.GetRegistrationStatusResponse, error) {
				return nil, status.Error(codes.Internal, "boom")
			},
		})
		go func() { _ = server.Serve(lis) }()
		defer server.Stop()

		cl, err := NewClient(RegistrationSagaConfig{Address: lis.Addr().String()}, logger)
		Expect(err).NotTo(HaveOccurred())
		defer func() { Expect(cl.Close()).To(Succeed()) }()

		_, err = cl.GetRegistrationStatus(context.Background(), "session-1", "client-1")
		Expect(err).To(MatchError(gateway.ErrFailedToCompleteRegistration))
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

	It("returns success from verify-email and completion lookups", func() {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			Skip("gRPC sockets are unavailable in this environment")
		}

		server := grpc.NewServer()
		registrationv1.RegisterRegistrationServiceServer(server, &fakeRegistrationServer{
			verifyFn: func(context.Context, *registrationv1.VerifyEmailRequest) (*registrationv1.VerifyEmailResponse, error) {
				return &registrationv1.VerifyEmailResponse{SessionId: "session-1", ClientId: "client-1", Status: "verifying_email"}, nil
			},
			statusFn: func(context.Context, *registrationv1.GetRegistrationStatusRequest) (*registrationv1.GetRegistrationStatusResponse, error) {
				return &registrationv1.GetRegistrationStatusResponse{SessionId: "session-1", ClientId: "client-1", UserId: "user-1", Status: "completed"}, nil
			},
		})
		go func() { _ = server.Serve(lis) }()
		defer server.Stop()

		cl, err := NewClient(RegistrationSagaConfig{Address: lis.Addr().String()}, logger)
		Expect(err).NotTo(HaveOccurred())
		defer func() { Expect(cl.Close()).To(Succeed()) }()

		result, err := cl.VerifyEmail(context.Background(), VerifyEmailRequest{
			SessionID: "session-1",
			ClientID:  "client-1",
			Code:      "123456",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&gateway.VerifyEmailResult{SessionID: "session-1", ClientID: "client-1", Status: "verifying_email"}))

		statusRes, err := cl.GetRegistrationStatus(context.Background(), "session-1", "client-1")
		Expect(err).NotTo(HaveOccurred())
		Expect(statusRes).To(Equal(&gateway.RegistrationStatus{SessionID: "session-1", ClientID: "client-1", UserID: "user-1", Status: "completed"}))
	})

	It("maps auth token issuance failures", func() {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			Skip("gRPC sockets are unavailable in this environment")
		}

		server := grpc.NewServer()
		authv1.RegisterAuthQueryServiceServer(server, &authTokenServer{
			issueFn: func(context.Context, *authv1.IssueRegistrationTokensRequest) (*authv1.IssueRegistrationTokensResponse, error) {
				return nil, status.Error(codes.Internal, "boom")
			},
		})
		go func() { _ = server.Serve(lis) }()
		defer server.Stop()

		cl, err := NewAuthClient(AuthServiceConfig{Address: lis.Addr().String()}, logger)
		Expect(err).NotTo(HaveOccurred())
		defer func() { Expect(cl.Close()).To(Succeed()) }()

		_, err = cl.IssueRegistrationTokens(context.Background(), "user-1")
		Expect(err).To(HaveOccurred())
	})

	It("issues registration tokens from auth grpc", func() {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			Skip("gRPC sockets are unavailable in this environment")
		}

		server := grpc.NewServer()
		authv1.RegisterAuthQueryServiceServer(server, authTokenServer{})
		go func() { _ = server.Serve(lis) }()
		defer server.Stop()

		cl, err := NewAuthClient(AuthServiceConfig{Address: lis.Addr().String()}, logger)
		Expect(err).NotTo(HaveOccurred())
		defer func() { Expect(cl.Close()).To(Succeed()) }()

		result, err := cl.IssueRegistrationTokens(context.Background(), "user-1")
		Expect(err).NotTo(HaveOccurred())
		Expect(result.UserID).To(Equal("user-1"))
		Expect(result.TokenType).To(Equal("Bearer"))
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
	if err != nil {
		Skip("gRPC sockets are unavailable in this environment")
		return nil, ""
	}

	server := grpc.NewServer()
	registrationv1.RegisterRegistrationServiceServer(server, &fakeRegistrationServer{startFn: startFn})

	go func() {
		_ = server.Serve(lis)
	}()

	return server, lis.Addr().String()
}

type authTokenServer struct {
	authv1.UnimplementedAuthQueryServiceServer
	issueFn func(context.Context, *authv1.IssueRegistrationTokensRequest) (*authv1.IssueRegistrationTokensResponse, error)
}

func (s authTokenServer) IssueRegistrationTokens(ctx context.Context, req *authv1.IssueRegistrationTokensRequest) (*authv1.IssueRegistrationTokensResponse, error) {
	if s.issueFn != nil {
		return s.issueFn(ctx, req)
	}
	return &authv1.IssueRegistrationTokensResponse{UserId: "user-1", TokenType: "Bearer"}, nil
}
