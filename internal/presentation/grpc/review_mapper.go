package grpc

import (
	"strings"

	gateway "api-gateway/internal/domain"
	reviewv1 "github.com/ofm-microservices/ofm-common/proto/review/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type reviewMapper struct{}

func newReviewMapper() ReviewMapper { return &reviewMapper{} }

func (m *reviewMapper) ToCreateReviewRequest(req gateway.CreateReviewRequest, buyerID string) *reviewv1.CreateReviewRequest {
	return &reviewv1.CreateReviewRequest{
		OrderId:     strings.TrimSpace(req.OrderID),
		BuyerUserId: strings.TrimSpace(buyerID),
		Content:     strings.TrimSpace(req.Content),
		Rating:      req.Rating,
		RequestedAt: strings.TrimSpace(req.RequestedAt),
	}
}

func (m *reviewMapper) ToCreateReviewResponse(res *reviewv1.CreateReviewResponse) *gateway.CreateReviewResult {
	if res == nil || res.GetReview() == nil {
		return &gateway.CreateReviewResult{}
	}
	review := res.GetReview()
	return &gateway.CreateReviewResult{
		ReviewID:  review.GetReviewId(),
		OrderID:   review.GetOrderId(),
		GigID:     review.GetGigId(),
		BuyerID:   review.GetBuyerUserId(),
		Content:   review.GetContent(),
		Rating:    review.GetRating(),
		CreatedAt: review.GetCreatedAt(),
	}
}

func (m *reviewMapper) ToError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.InvalidArgument:
		return gateway.ErrInvalidReviewContent
	case codes.FailedPrecondition:
		return gateway.ErrOrderNotAcceptable
	default:
		return err
	}
}
