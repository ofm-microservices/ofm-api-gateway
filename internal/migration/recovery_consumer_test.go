package migration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ofm-microservices/ofm-common/pkg/migration/events"
)

func TestRecoveryReplayClassifiesTransientAndPermanentResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/transient" {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusConflict)
	}))
	defer server.Close()

	consumer := &recoveryConsumer{base: server.URL, client: server.Client(), maxAttempts: 5}
	transient := consumer.replay(context.Background(), events.Envelope{CommandID: "c1", CommandMethod: http.MethodPost, CommandPath: "/transient"})
	if transient == nil || !transient.transient {
		t.Fatalf("expected 503 to be transient, got %#v", transient)
	}
	permanent := consumer.replay(context.Background(), events.Envelope{CommandID: "c2", CommandMethod: http.MethodPost, CommandPath: "/permanent"})
	if permanent == nil || permanent.transient {
		t.Fatalf("expected 409 to be permanent, got %#v", permanent)
	}
}

func TestSplitBrokers(t *testing.T) {
	got := splitBrokers("kafka-a:9092, kafka-b:9092\nkafka-c:9092")
	if len(got) != 3 || got[0] != "kafka-a:9092" || got[2] != "kafka-c:9092" {
		t.Fatalf("unexpected broker split: %#v", got)
	}
}
