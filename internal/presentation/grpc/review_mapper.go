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
		OrderId:       strings.TrimSpace(req.OrderID),
		BuyerUserId:   strings.TrimSpace(buyerID),
		BuyerUsername: strings.TrimSpace(req.BuyerUsername),
		Content:       strings.TrimSpace(req.Content),
		Rating:        req.Rating,
		RequestedAt:   strings.TrimSpace(req.RequestedAt),
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

func (m *reviewMapper) ToGetGigReviewsRequest(req gateway.GetGigReviewsRequest) *reviewv1.ListGigReviewsRequest {
	return &reviewv1.ListGigReviewsRequest{
		GigId:  strings.TrimSpace(req.GigID),
		Cursor: strings.TrimSpace(req.Cursor),
	}
}

func (m *reviewMapper) ToGetGigReviewsResponse(res *reviewv1.ListGigReviewsResponse) *gateway.GetGigReviewsResult {
	if res == nil {
		return &gateway.GetGigReviewsResult{}
	}
	out := &gateway.GetGigReviewsResult{
		Reviews: &gateway.ReviewList{
			Cursor:  res.GetCursor(),
			HasMore: res.GetHasMore(),
		},
	}
	if len(res.GetReviews()) > 0 {
		out.Reviews.Items = make([]gateway.Review, 0, len(res.GetReviews()))
		for _, item := range res.GetReviews() {
			out.Reviews.Items = append(out.Reviews.Items, gateway.Review{
				ReviewID:    item.GetReviewId(),
				OrderID:     item.GetOrderId(),
				GigID:       item.GetGigId(),
				BuyerUserID: item.GetBuyerUserId(),
				Content:     item.GetContent(),
				Rating:      item.GetRating(),
				CreatedAt:   item.GetCreatedAt(),
				Author:      toReviewAuthor(item.GetAuthor()),
			})
		}
	}
	return out
}

func (m *reviewMapper) ToGetReviewsBySellerUsernameRequest(req gateway.GetReviewsBySellerUsernameRequest) *reviewv1.GetReviewsBySellerUsernameRequest {
	return &reviewv1.GetReviewsBySellerUsernameRequest{
		Username: strings.TrimSpace(req.Username),
		Cursor:   strings.TrimSpace(req.Cursor),
	}
}

func (m *reviewMapper) ToGetReviewsBySellerUsernameResponse(res *reviewv1.ListSellerReviewsResponse) *gateway.GetReviewsBySellerUsernameResult {
	if res == nil {
		return &gateway.GetReviewsBySellerUsernameResult{}
	}
	out := &gateway.GetReviewsBySellerUsernameResult{
		Reviews: &gateway.ReviewList{
			Cursor:  res.GetCursor(),
			HasMore: res.GetHasMore(),
		},
	}
	if len(res.GetReviews()) > 0 {
		out.Reviews.Items = make([]gateway.Review, 0, len(res.GetReviews()))
		for _, item := range res.GetReviews() {
			out.Reviews.Items = append(out.Reviews.Items, gateway.Review{
				ReviewID:       item.GetReviewId(),
				OrderID:        item.GetOrderId(),
				GigID:          item.GetGigId(),
				BuyerUserID:    item.GetBuyerUserId(),
				Content:        item.GetContent(),
				Rating:         item.GetRating(),
				CreatedAt:      item.GetCreatedAt(),
				SellerUsername: item.GetSellerUsername(),
				Author:         toReviewAuthor(item.GetAuthor()),
			})
		}
	}
	return out
}

func (m *reviewMapper) ToGetGigReviewsSummaryRequest(req gateway.GetGigReviewsSummaryRequest) *reviewv1.GetGigRatingSummaryRequest {
	return &reviewv1.GetGigRatingSummaryRequest{GigId: strings.TrimSpace(req.GigID)}
}

func (m *reviewMapper) ToGetGigReviewsSummaryResponse(res *reviewv1.RatingSummary) *gateway.ReviewSummary {
	if res == nil {
		return nil
	}
	return &gateway.ReviewSummary{
		RatingAvg:    res.GetRatingAvg(),
		TotalReviews: res.GetTotalReviews(),
		Stars5:       res.GetStars_5(),
		Stars4:       res.GetStars_4(),
		Stars3:       res.GetStars_3(),
		Stars2:       res.GetStars_2(),
		Stars1:       res.GetStars_1(),
	}
}

func (m *reviewMapper) ToGetUserRatingSummaryByUsernameRequest(req gateway.GetUserRatingSummaryByUsernameRequest) *reviewv1.GetUserRatingSummaryByUsernameRequest {
	return &reviewv1.GetUserRatingSummaryByUsernameRequest{Username: strings.TrimSpace(req.Username)}
}

func (m *reviewMapper) ToGetUserRatingSummaryByUsernameResponse(res *reviewv1.RatingSummary) *gateway.ReviewSummary {
	if res == nil {
		return nil
	}
	return &gateway.ReviewSummary{
		RatingAvg:    res.GetRatingAvg(),
		TotalReviews: res.GetTotalReviews(),
		Stars5:       res.GetStars_5(),
		Stars4:       res.GetStars_4(),
		Stars3:       res.GetStars_3(),
		Stars2:       res.GetStars_2(),
		Stars1:       res.GetStars_1(),
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
		if st.Message() == gateway.ErrInvalidSellerID.Error() {
			return gateway.ErrInvalidSellerID
		}
		return gateway.ErrInvalidReviewContent
	case codes.FailedPrecondition:
		return gateway.ErrOrderNotAcceptable
	default:
		return err
	}
}

func toReviewAuthor(res *reviewv1.ReviewAuthor) *gateway.ReviewAuthor {
	if res == nil {
		return nil
	}
	return &gateway.ReviewAuthor{
		UserID:      res.GetUserId(),
		Username:    res.GetUsername(),
		DisplayName: res.GetDisplayName(),
		AvatarURL:   res.GetAvatarUrl(),
	}
}
