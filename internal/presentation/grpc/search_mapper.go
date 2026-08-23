package grpc

import (
	gateway "api-gateway/internal/domain"
	searchv1 "github.com/ofm-microservices/ofm-common/proto/search/v1"
)

type searchMapper struct{}

func newSearchMapper() SearchMapper { return &searchMapper{} }

func (m *searchMapper) ToSearchRequest(req gateway.SearchRequest) *searchv1.SearchRequest {
	return &searchv1.SearchRequest{
		Query:  req.Query,
		Sort:   req.Sort,
		Order:  req.Order,
		Cursor: req.Cursor,
	}
}

func (m *searchMapper) ToSearchResponse(res *searchv1.SearchResponse) *gateway.SearchResponse {
	if res == nil {
		return &gateway.SearchResponse{}
	}
	out := &gateway.SearchResponse{
		Cursor:  res.GetCursor(),
		HasMore: res.GetHasMore(),
	}
	out.Items = make([]gateway.SearchResult, 0, len(res.GetServices()))
	for _, item := range res.GetServices() {
		out.Items = append(out.Items, gateway.SearchResult{
			ID:             item.GetId(),
			Title:          item.GetTitle(),
			Description:    item.GetDescription(),
			Picture:        item.GetPicture(),
			ReviewsCount:   item.GetReviewsCount(),
			Rating:         item.GetRating(),
			MinPrice:       item.GetMinPrice(),
			Slug:           item.GetSlug(),
			FreelancerID:   item.GetFreelancerId(),
			SellerUsername: item.GetSellerUsername(),
			PublishedAt:    item.GetPublishedAt(),
		})
	}
	return out
}
