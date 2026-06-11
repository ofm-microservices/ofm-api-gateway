package http

import gateway "api-gateway/internal/domain"

var (
	_ = gateway.GetOrderPreviewByIDResult{}
	_ = gateway.GetOrderRequirementsByIDResult{}
	_ = gateway.GetOrderDeliveryByIDResult{}
)

// swaggerGetOrderPreviewByIDDoc documents the authenticated order preview endpoint.
//
// @Summary Get order preview
// @Description Loads the authenticated user-scoped order preview payload.
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Param role query string false "Participant role"
// @Success 200 {object} gateway.GetOrderPreviewByIDResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id} [get]
func swaggerGetOrderPreviewByIDDoc() {}

// swaggerGetOrderRequirementsByIDDoc documents the authenticated requirements endpoint.
//
// @Summary Get order requirements
// @Description Loads the buyer requirements page for an authenticated user-scoped order.
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Success 200 {object} gateway.GetOrderRequirementsByIDResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/requirements [get]
func swaggerGetOrderRequirementsByIDDoc() {}

// swaggerGetOrderDeliveryByIDDoc documents the authenticated delivery endpoint.
//
// @Summary Get order delivery
// @Description Loads the seller delivery page for an authenticated user-scoped order.
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param order_id path string true "Order ID"
// @Success 200 {object} gateway.GetOrderDeliveryByIDResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/orders/{order_id}/delivery [get]
func swaggerGetOrderDeliveryByIDDoc() {}
