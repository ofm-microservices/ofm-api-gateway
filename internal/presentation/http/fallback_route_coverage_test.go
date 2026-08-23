package http

import "testing"

// TestFallbackWritePathCoverage keeps the migration boundary aligned with the
// complete v2 write surface. The middleware is intentionally global, so every
// listed route must remain inside the versioned write namespace.
func TestFallbackWritePathCoverage(t *testing.T) {
	routes := []struct {
		method string
		path   string
		owner  string
	}{
		{"POST", "/api/v2/auth/sign-up", "auth/registration"},
		{"POST", "/api/v2/auth/sign-in", "auth"},
		{"POST", "/api/v2/auth/refresh", "auth"},
		{"POST", "/api/v2/auth/sign-out", "auth"},
		{"POST", "/api/v2/auth/sign-up/verify-email", "registration"},
		{"POST", "/api/v2/auth/sign-up/complete", "registration"},
		{"POST", "/api/v2/freelancer/onboarding/start", "payment"},
		{"POST", "/api/v2/gigs/drafts", "gig"},
		{"PATCH", "/api/v2/gigs/{gig_id}/basic-info", "gig"},
		{"PUT", "/api/v2/gigs/{gig_id}/packages", "gig"},
		{"PUT", "/api/v2/gigs/{gig_id}/requirements", "gig"},
		{"PUT", "/api/v2/gigs/{gig_id}/media", "gig"},
		{"POST", "/api/v2/gigs/{gig_id}/publish", "gig"},
		{"POST", "/api/v2/orders/start", "order-saga"},
		{"POST", "/api/v2/orders/{order_id}/confirm", "order-saga/payment"},
		{"POST", "/api/v2/orders/{order_id}/requirements", "order-saga"},
		{"POST", "/api/v2/orders/{order_id}/message", "order-saga/chat"},
		{"POST", "/api/v2/orders/{order_id}/deliver", "order-saga"},
		{"POST", "/api/v2/orders/{order_id}/accept", "order-saga"},
		{"POST", "/api/v2/orders/{order_id}/request-revision", "order-saga"},
		{"POST", "/api/v2/orders/{order_id}/dispute", "order-saga"},
		{"POST", "/api/v2/orders/{order_id}/dispute/resolve", "order-saga"},
		{"POST", "/api/v2/orders/{order_id}/reviews", "review"},
		{"POST", "/api/v2/users/{username}/orders/{order_id}/chat/messages", "chat"},
		{"PATCH", "/api/v2/users/{username}/orders/{order_id}/chat/messages/{message_id}", "chat"},
		{"DELETE", "/api/v2/users/{username}/orders/{order_id}/chat/messages/{message_id}", "chat"},
		{"POST", "/api/v2/users/{username}/orders/{order_id}/chat/attachments/upload-url", "chat/file"},
		{"POST", "/api/v2/users/{username}/orders/{order_id}/chat/attachments/complete", "chat/file"},
	}

	for _, route := range routes {
		if !isFallbackWritePath(route.path) {
			t.Errorf("%s %s (%s) is outside the fallback namespace", route.method, route.path, route.owner)
		}
		if !isWriteMethod(route.method) {
			t.Errorf("%s %s (%s) is not classified as a write", route.method, route.path, route.owner)
		}
	}
}
