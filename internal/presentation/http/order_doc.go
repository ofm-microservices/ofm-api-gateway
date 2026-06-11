package http

import gateway "api-gateway/internal/domain"

var (
	_ = gateway.CreateOrderRequest{}
	_ = gateway.CreateOrderResult{}
	_ = gateway.ConfirmOrderRequest{}
	_ = gateway.ConfirmOrderResult{}
	_ = gateway.SubmitOrderRequirementsRequest{}
	_ = gateway.SubmitOrderRequirementsResult{}
	_ = gateway.SubmitOrderMessageRequest{}
	_ = gateway.SubmitOrderMessageResult{}
	_ = gateway.DeliverOrderRequest{}
	_ = gateway.DeliverOrderResult{}
	_ = gateway.AcceptDeliveryRequest{}
	_ = gateway.AcceptDeliveryResult{}
	_ = gateway.RequestRevisionRequest{}
	_ = gateway.RequestRevisionResult{}
	_ = gateway.OpenDisputeRequest{}
	_ = gateway.OpenDisputeResult{}
)

// swaggerCreateOrderDoc documents the public order start endpoint.
//
// @Summary Start order
// @Description Starts the public order flow for a gig package.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body gateway.CreateOrderRequest true "Order start payload"
// @Success 201 {object} gateway.CreateOrderResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/start [post]
func swaggerCreateOrderDoc() {}

// swaggerConfirmOrderDoc documents the payment confirmation endpoint.
//
// @Summary Confirm order
// @Description Confirms the order and creates the payment session.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.ConfirmOrderRequest true "Order confirmation payload"
// @Success 200 {object} gateway.ConfirmOrderResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/confirm [post]
func swaggerConfirmOrderDoc() {}

// swaggerSubmitRequirementsDoc documents the buyer requirements submission endpoint.
//
// @Summary Submit requirements
// @Description Submits the buyer answers for the order requirements step.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.SubmitOrderRequirementsRequest true "Requirements payload"
// @Success 200 {object} gateway.SubmitOrderRequirementsResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/requirements [post]
func swaggerSubmitRequirementsDoc() {}

// swaggerSubmitMessageDoc documents the buyer message submission endpoint.
//
// @Summary Submit message
// @Description Submits the buyer initial message for the order.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.SubmitOrderMessageRequest true "Message payload"
// @Success 200 {object} gateway.SubmitOrderMessageResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/message [post]
func swaggerSubmitMessageDoc() {}

// swaggerDeliverOrderDoc documents the seller delivery endpoint.
//
// @Summary Deliver order
// @Description Submits the seller delivery payload.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.DeliverOrderRequest true "Delivery payload"
// @Success 200 {object} gateway.DeliverOrderResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/deliver [post]
func swaggerDeliverOrderDoc() {}

// swaggerAcceptDeliveryDoc documents the buyer acceptance endpoint.
//
// @Summary Accept delivery
// @Description Confirms the buyer accepted the seller delivery.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.AcceptDeliveryRequest true "Acceptance payload"
// @Success 200 {object} gateway.AcceptDeliveryResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/accept [post]
func swaggerAcceptDeliveryDoc() {}

// swaggerRequestRevisionDoc documents the revision request endpoint.
//
// @Summary Request revision
// @Description Asks the seller for revisions on the current delivery.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.RequestRevisionRequest true "Revision payload"
// @Success 200 {object} gateway.RequestRevisionResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/request-revision [post]
func swaggerRequestRevisionDoc() {}

// swaggerOpenDisputeDoc documents the buyer dispute endpoint.
//
// @Summary Open dispute
// @Description Opens a buyer dispute for the current delivery.
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param request body gateway.OpenDisputeRequest true "Dispute payload"
// @Success 200 {object} gateway.OpenDisputeResult
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /orders/{order_id}/dispute [post]
func swaggerOpenDisputeDoc() {}
