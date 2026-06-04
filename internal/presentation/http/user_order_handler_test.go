package http

import (
	gateway "api-gateway/internal/domain"
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"github.com/valyala/fasthttp"
)

type userOrderServiceStub struct{}

func (userOrderServiceStub) GetOrderPreviewByID(_ context.Context, _ gateway.GetOrderPreviewByIDRequest) (*gateway.GetOrderPreviewByIDResult, error) {
	return &gateway.GetOrderPreviewByIDResult{}, nil
}

func (userOrderServiceStub) GetOrderRequirementsByID(_ context.Context, _ gateway.GetOrderRequirementsByIDRequest) (*gateway.GetOrderRequirementsByIDResult, error) {
	return &gateway.GetOrderRequirementsByIDResult{}, nil
}

func (userOrderServiceStub) GetOrderDeliveryByID(_ context.Context, _ gateway.GetOrderDeliveryByIDRequest) (*gateway.GetOrderDeliveryByIDResult, error) {
	return &gateway.GetOrderDeliveryByIDResult{}, nil
}

func TestUserOrderHandlerMapsOwnershipFailuresToForbidden(t *testing.T) {
	t.Helper()
	logger, err := logging.New("api-gateway", "test", "debug")
	if err != nil {
		t.Fatalf("logging.New: %v", err)
	}
	handler, err := NewUserOrderHandler(userOrderServiceStub{}, "secret", logger)
	if err != nil {
		t.Fatalf("NewUserOrderHandler: %v", err)
	}
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	if err := handler.(*userOrderHandler).mapError(ctx, gateway.ErrOrderNotOwned); err != nil {
		t.Fatalf("mapError: %v", err)
	}
	if got := ctx.Response().StatusCode(); got != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, got)
	}
}

func TestUserOrderHandlerMapsDeliveryNotFoundToNotFound(t *testing.T) {
	t.Helper()
	logger, err := logging.New("api-gateway", "test", "debug")
	if err != nil {
		t.Fatalf("logging.New: %v", err)
	}
	handler, err := NewUserOrderHandler(userOrderServiceStub{}, "secret", logger)
	if err != nil {
		t.Fatalf("NewUserOrderHandler: %v", err)
	}
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	if err := handler.(*userOrderHandler).mapError(ctx, gateway.ErrOrderDeliveryNotFound); err != nil {
		t.Fatalf("mapError: %v", err)
	}
	if got := ctx.Response().StatusCode(); got != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, got)
	}
}
