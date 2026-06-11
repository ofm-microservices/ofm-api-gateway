package http

import gateway "api-gateway/internal/domain"

var _ = gateway.CreateReviewResult{}

// swaggerCreateReviewDoc documents the buyer review creation endpoint.
//
// @Summary Create review
// @Description Submits a buyer-authenticated review for a completed order.
// @Tags reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.CreateReviewRequest true "Review payload"
// @Success 201 {object} gateway.CreateReviewResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/reviews [post]
func swaggerCreateReviewDoc() {}
