package nats

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"api-gateway/config"
	gateway "api-gateway/internal/domain"
	"github.com/nats-io/nats.go"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNATS(t *testing.T) {
	t.Helper()
	RegisterFailHandler(Fail)
	RunSpecs(t, "NATS Suite")
}

var (
	natsContainer     testcontainers.Container
	natsAuthContainer testcontainers.Container
	natsCfg           config.NATSConfig
	natsAuthCfg       config.NATSConfig
	logger            logging.Logger
)

var _ = BeforeSuite(func() {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	natsContainer, natsCfg = startNATSContainer(ctx, "", "")
	natsAuthContainer, natsAuthCfg = startNATSContainer(ctx, "ofm", "secret")

	var err error
	logger, err = logging.New("api-gateway", "test", "debug")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if natsContainer != nil {
		Expect(natsContainer.Terminate(context.Background())).To(Succeed())
	}
	if natsAuthContainer != nil {
		Expect(natsAuthContainer.Terminate(context.Background())).To(Succeed())
	}
})

var _ = Describe("publisher", func() {
	It("validates constructor input", func() {
		pub, err := NewPublisher(config.NATSConfig{}, logger)
		Expect(pub).To(BeNil())
		Expect(err).To(MatchError(ErrEmptyNATSURL))

		pub, err = NewPublisher(config.NATSConfig{URL: "nats://127.0.0.1:4222"}, logger)
		Expect(pub).To(BeNil())
		Expect(err).To(MatchError(ErrEmptySubject))

		pub, err = NewPublisher(config.NATSConfig{
			URL:                   "nats://127.0.0.1:4222",
			SagaCreateAuthSubject: "saga.auth.create",
		}, nil)
		Expect(pub).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("publishes signup requests to a real nats server", func() {
		nc, err := nats.Connect(natsCfg.URL)
		Expect(err).NotTo(HaveOccurred())
		defer nc.Close()

		received := make(chan gateway.SignUpRequest, 1)
		_, err = nc.Subscribe(natsCfg.SagaCreateAuthSubject, func(msg *nats.Msg) {
			var req gateway.SignUpRequest
			if json.Unmarshal(msg.Data, &req) == nil {
				received <- req
			}
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(nc.Flush()).To(Succeed())

		pubAny, err := NewPublisher(natsCfg, logger)
		Expect(err).NotTo(HaveOccurred())
		pub := pubAny.(*publisher)
		defer pub.Close()

		result, err := pub.StartRegistration(context.Background(), gateway.SignUpRequest{
			Email:     "alex@example.com",
			Username:  "alex",
			Password:  "password123",
			FirstName: "Alex",
			Surname:   "Doe",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(&gateway.SignUpResult{Status: "pending"}))
		Eventually(received, time.Second).Should(Receive(Equal(gateway.SignUpRequest{
			Email:     "alex@example.com",
			Username:  "alex",
			Password:  "password123",
			FirstName: "Alex",
			Surname:   "Doe",
		})))
	})

	It("uses the caller deadline when flushing", func() {
		nc, err := nats.Connect(natsCfg.URL)
		Expect(err).NotTo(HaveOccurred())
		defer nc.Close()

		pubAny, err := NewPublisher(natsCfg, logger)
		Expect(err).NotTo(HaveOccurred())
		pub := pubAny.(*publisher)
		defer pub.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		result, err := pub.StartRegistration(ctx, gateway.SignUpRequest{
			Email:    "alex@example.com",
			Username: "alex",
			Password: "password123",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Status).To(Equal("pending"))
	})

	It("supports authenticated nats connections", func() {
		pub, err := NewPublisher(natsAuthCfg, logger)
		Expect(err).NotTo(HaveOccurred())
		defer pub.Close()
	})

	It("wraps publish failures when the connection is closed", func() {
		pubAny, err := NewPublisher(natsCfg, logger)
		Expect(err).NotTo(HaveOccurred())

		pub := pubAny.(*publisher)
		pub.nc.Close()

		result, err := pub.StartRegistration(context.Background(), gateway.SignUpRequest{
			Email:    "alex@example.com",
			Username: "alex",
			Password: "password123",
		})

		Expect(result).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("publish to nats"))
	})

	It("covers wrapper helpers", func() {
		Expect(WrapConnectToNATSError(context.Canceled)).To(MatchError(ContainSubstring("connect to nats")))
		Expect(WrapMarshalEventError(context.Canceled)).To(MatchError(ContainSubstring("marshal event")))
		Expect(WrapPublishToNATSError("subject", context.Canceled)).To(MatchError(ContainSubstring("publish to nats (subject)")))
		Expect(WrapFlushNATSError(context.Canceled)).To(MatchError(ContainSubstring("flush nats publisher")))
	})

	It("closes a nil publisher safely", func() {
		var pub *publisher
		Expect(func() { pub.Close() }).NotTo(Panic())
	})

	It("covers direct wrapper helpers", func() {
		Expect(WrapConnectToNATSError(context.Canceled).Error()).To(ContainSubstring("connect to nats"))
		Expect(WrapMarshalEventError(context.Canceled).Error()).To(ContainSubstring("marshal event"))
		Expect(WrapPublishToNATSError("subject", context.Canceled).Error()).To(ContainSubstring("publish to nats (subject)"))
		Expect(WrapFlushNATSError(context.Canceled).Error()).To(ContainSubstring("flush nats publisher"))
	})
})

func startNATSContainer(ctx context.Context, user, password string) (testcontainers.Container, config.NATSConfig) {
	defer func() {
		if r := recover(); r != nil {
			Skip("docker-based integration tests are unavailable in this environment")
		}
	}()

	cmd := []string{}
	if user != "" {
		cmd = []string{"--user", user, "--pass", password}
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "nats:2.11-alpine",
			ExposedPorts: []string{"4222/tcp"},
			Cmd:          cmd,
			WaitingFor:   wait.ForListeningPort("4222/tcp"),
		},
		Started: true,
	})
	if err != nil {
		Skip("docker-based integration tests are unavailable in this environment")
		return nil, config.NATSConfig{}
	}

	host, err := container.Host(ctx)
	Expect(err).NotTo(HaveOccurred())
	port, err := container.MappedPort(ctx, "4222/tcp")
	Expect(err).NotTo(HaveOccurred())

	return container, config.NATSConfig{
		URL:                   "nats://" + host + ":" + port.Port(),
		User:                  user,
		Password:              password,
		SagaCreateAuthSubject: "saga.auth.create",
	}
}
