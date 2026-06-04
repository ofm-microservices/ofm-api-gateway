package grpc

import (
	"strings"

	gateway "api-gateway/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ordercheckoutv1 "github.com/ofm-microservices/ofm-common/proto/ordercheckout/v1"
)

type orderCheckoutMapper struct{}

func newOrderCheckoutMapper() OrderCheckoutMapper { return &orderCheckoutMapper{} }

func (m *orderCheckoutMapper) ToStartOrderRequest(req gateway.CreateOrderRequest) *ordercheckoutv1.StartOrderRequest {
	return &ordercheckoutv1.StartOrderRequest{
		BuyerUserId:          strings.TrimSpace(req.BuyerID),
		BuyerEmail:           strings.TrimSpace(req.BuyerEmail),
		GigId:                strings.TrimSpace(req.GigID),
		PackageId:            strings.TrimSpace(req.PackageID),
		IdempotencyKey:       strings.TrimSpace(req.IdempotencyKey),
		CorrelationId:        strings.TrimSpace(req.RealtimeConnectionID),
		RealtimeConnectionId: strings.TrimSpace(req.RealtimeConnectionID),
		RequestedAt:          strings.TrimSpace(req.RequestedAt),
	}
}

func (m *orderCheckoutMapper) ToStartOrderResponse(res *ordercheckoutv1.StartOrderResponse) *gateway.CreateOrderResult {
	if res == nil {
		return &gateway.CreateOrderResult{}
	}
	out := &gateway.CreateOrderResult{
		SagaID:  res.GetSagaId(),
		OrderID: res.GetOrderId(),
		Status:  res.GetStatus(),
	}
	if snap := res.GetSnapshot(); snap != nil {
		out.Snapshot = &gateway.OrderSnapshot{
			GigID:              snap.GetGigId(),
			PackageID:          snap.GetPackageId(),
			SellerUserID:       snap.GetSellerUserId(),
			GigTitle:           snap.GetGigTitle(),
			PackageTitle:       snap.GetPackageTitle(),
			PackageDescription: snap.GetPackageDescription(),
			PriceCents:         snap.GetPriceAmount(),
			Currency:           snap.GetPriceCurrency(),
			DeliveryDays:       snap.GetDeliveryDays(),
			RevisionCount:      snap.GetRevisionCount(),
		}
	}
	for _, q := range res.GetQuestions() {
		out.Questions = append(out.Questions, gateway.OrderQuestion{
			ID:        q.GetQuestionId(),
			Text:      q.GetText(),
			Type:      q.GetType(),
			Required:  q.GetRequired(),
			SortOrder: q.GetSortOrder(),
		})
	}
	return out
}

func (m *orderCheckoutMapper) ToConfirmOrderRequest(req gateway.ConfirmOrderRequest) *ordercheckoutv1.ConfirmOrderRequest {
	return &ordercheckoutv1.ConfirmOrderRequest{
		OrderId:     strings.TrimSpace(req.OrderID),
		BuyerUserId: strings.TrimSpace(req.BuyerID),
		RequestedAt: strings.TrimSpace(req.RequestedAt),
	}
}

func (m *orderCheckoutMapper) ToConfirmOrderResponse(res *ordercheckoutv1.ConfirmOrderResponse) *gateway.ConfirmOrderResult {
	if res == nil {
		return &gateway.ConfirmOrderResult{}
	}
	return &gateway.ConfirmOrderResult{
		OrderID:     res.GetOrderId(),
		Status:      res.GetStatus(),
		CheckoutURL: res.GetCheckoutUrl(),
		PaymentID:   res.GetPaymentId(),
	}
}

func (m *orderCheckoutMapper) ToSubmitRequirementsRequest(req gateway.SubmitOrderRequirementsRequest) *ordercheckoutv1.SubmitRequirementsRequest {
	answers := make([]*ordercheckoutv1.OrderRequirementAnswer, 0, len(req.Answers))
	for _, a := range req.Answers {
		answers = append(answers, &ordercheckoutv1.OrderRequirementAnswer{QuestionId: a.QuestionID, Value: a.Value})
	}
	return &ordercheckoutv1.SubmitRequirementsRequest{
		OrderId:     strings.TrimSpace(req.OrderID),
		BuyerUserId: strings.TrimSpace(req.BuyerID),
		Answers:     answers,
		RequestedAt: strings.TrimSpace(req.RequestedAt),
	}
}

func (m *orderCheckoutMapper) ToSubmitRequirementsResponse(res *ordercheckoutv1.SubmitRequirementsResponse) *gateway.SubmitOrderRequirementsResult {
	if res == nil {
		return &gateway.SubmitOrderRequirementsResult{}
	}
	return &gateway.SubmitOrderRequirementsResult{OrderID: res.GetOrderId(), Status: res.GetStatus(), CurrentStep: res.GetCurrentStep()}
}

func (m *orderCheckoutMapper) ToSubmitMessageRequest(req gateway.SubmitOrderMessageRequest) *ordercheckoutv1.SubmitMessageRequest {
	return &ordercheckoutv1.SubmitMessageRequest{OrderId: strings.TrimSpace(req.OrderID), BuyerUserId: strings.TrimSpace(req.BuyerID), Message: strings.TrimSpace(req.Message), RequestedAt: strings.TrimSpace(req.RequestedAt)}
}

func (m *orderCheckoutMapper) ToSubmitMessageResponse(res *ordercheckoutv1.SubmitMessageResponse) *gateway.SubmitOrderMessageResult {
	if res == nil {
		return &gateway.SubmitOrderMessageResult{}
	}
	return &gateway.SubmitOrderMessageResult{OrderID: res.GetOrderId(), Status: res.GetStatus(), CurrentStep: res.GetCurrentStep()}
}

func (m *orderCheckoutMapper) ToCreateAttachmentUploadURLRequest(req gateway.CreateOrderAttachmentUploadURLRequest) *ordercheckoutv1.CreateAttachmentUploadURLRequest {
	return &ordercheckoutv1.CreateAttachmentUploadURLRequest{OrderId: strings.TrimSpace(req.OrderID), BuyerUserId: strings.TrimSpace(req.BuyerID), FileName: strings.TrimSpace(req.FileName), MimeType: strings.TrimSpace(req.MimeType), SizeBytes: req.SizeBytes, RequestedAt: strings.TrimSpace(req.RequestedAt)}
}

func (m *orderCheckoutMapper) ToCreateAttachmentUploadURLResponse(res *ordercheckoutv1.CreateAttachmentUploadURLResponse) *gateway.CreateOrderAttachmentUploadURLResult {
	if res == nil {
		return &gateway.CreateOrderAttachmentUploadURLResult{}
	}
	return &gateway.CreateOrderAttachmentUploadURLResult{OrderID: res.GetOrderId(), AttachmentID: res.GetAttachmentId(), UploadURL: res.GetUploadUrl(), FileKey: res.GetFileKey()}
}

func (m *orderCheckoutMapper) ToCompleteAttachmentUploadRequest(req gateway.CompleteOrderAttachmentUploadRequest) *ordercheckoutv1.CompleteAttachmentUploadRequest {
	return &ordercheckoutv1.CompleteAttachmentUploadRequest{OrderId: strings.TrimSpace(req.OrderID), BuyerUserId: strings.TrimSpace(req.BuyerID), AttachmentId: strings.TrimSpace(req.AttachmentID), FileKey: strings.TrimSpace(req.FileKey), RequestedAt: strings.TrimSpace(req.RequestedAt)}
}

func (m *orderCheckoutMapper) ToCompleteAttachmentUploadResponse(res *ordercheckoutv1.CompleteAttachmentUploadResponse) *gateway.CompleteOrderAttachmentUploadResult {
	if res == nil {
		return &gateway.CompleteOrderAttachmentUploadResult{}
	}
	return &gateway.CompleteOrderAttachmentUploadResult{OrderID: res.GetOrderId(), AttachmentID: res.GetAttachmentId(), Status: res.GetStatus()}
}

func (m *orderCheckoutMapper) ToDeliverOrderRequest(req gateway.DeliverOrderRequest) *ordercheckoutv1.DeliverOrderRequest {
	return &ordercheckoutv1.DeliverOrderRequest{
		OrderId:         strings.TrimSpace(req.OrderID),
		SellerUserId:    strings.TrimSpace(req.SellerID),
		DeliveryMessage: strings.TrimSpace(req.DeliveryMessage),
		AttachmentIds:   append([]string(nil), req.AttachmentIDs...),
		RequestedAt:     strings.TrimSpace(req.RequestedAt),
	}
}

func (m *orderCheckoutMapper) ToDeliverOrderResponse(res *ordercheckoutv1.DeliverOrderResponse) *gateway.DeliverOrderResult {
	if res == nil {
		return &gateway.DeliverOrderResult{}
	}
	return &gateway.DeliverOrderResult{OrderID: res.GetOrderId(), Status: res.GetStatus(), CurrentStep: res.GetCurrentStep()}
}

func (m *orderCheckoutMapper) ToAcceptDeliveryRequest(req gateway.AcceptDeliveryRequest) *ordercheckoutv1.AcceptDeliveryRequest {
	return &ordercheckoutv1.AcceptDeliveryRequest{OrderId: strings.TrimSpace(req.OrderID), BuyerUserId: strings.TrimSpace(req.BuyerID), RequestedAt: strings.TrimSpace(req.RequestedAt)}
}

func (m *orderCheckoutMapper) ToAcceptDeliveryResponse(res *ordercheckoutv1.AcceptDeliveryResponse) *gateway.AcceptDeliveryResult {
	if res == nil {
		return &gateway.AcceptDeliveryResult{}
	}
	return &gateway.AcceptDeliveryResult{OrderID: res.GetOrderId(), Status: res.GetStatus(), CurrentStep: res.GetCurrentStep()}
}

func (m *orderCheckoutMapper) ToRequestRevisionRequest(req gateway.RequestRevisionRequest) *ordercheckoutv1.RequestRevisionRequest {
	return &ordercheckoutv1.RequestRevisionRequest{OrderId: strings.TrimSpace(req.OrderID), BuyerUserId: strings.TrimSpace(req.BuyerID), Reason: strings.TrimSpace(req.Reason), RequestedAt: strings.TrimSpace(req.RequestedAt)}
}

func (m *orderCheckoutMapper) ToRequestRevisionResponse(res *ordercheckoutv1.RequestRevisionResponse) *gateway.RequestRevisionResult {
	if res == nil {
		return &gateway.RequestRevisionResult{}
	}
	return &gateway.RequestRevisionResult{OrderID: res.GetOrderId(), Status: res.GetStatus(), CurrentStep: res.GetCurrentStep()}
}

func (m *orderCheckoutMapper) ToOpenDisputeRequest(req gateway.OpenDisputeRequest) *ordercheckoutv1.OpenDisputeRequest {
	return &ordercheckoutv1.OpenDisputeRequest{OrderId: strings.TrimSpace(req.OrderID), BuyerUserId: strings.TrimSpace(req.BuyerID), Reason: strings.TrimSpace(req.Reason), RequestedAt: strings.TrimSpace(req.RequestedAt)}
}

func (m *orderCheckoutMapper) ToOpenDisputeResponse(res *ordercheckoutv1.OpenDisputeResponse) *gateway.OpenDisputeResult {
	if res == nil {
		return &gateway.OpenDisputeResult{}
	}
	return &gateway.OpenDisputeResult{OrderID: res.GetOrderId(), Status: res.GetStatus(), CurrentStep: res.GetCurrentStep()}
}

func (m *orderCheckoutMapper) ToError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.InvalidArgument:
		return gateway.ErrFailedToCreateOrder
	case codes.PermissionDenied:
		if st.Message() == gateway.ErrOrderNotOwned.Error() {
			return gateway.ErrOrderNotOwned
		}
		return gateway.ErrOrderNotOwned
	case codes.NotFound:
		if st.Message() == gateway.ErrGigNotFound.Error() {
			return gateway.ErrGigNotFound
		}
		return gateway.ErrFailedToGetGig
	case codes.FailedPrecondition:
		switch st.Message() {
		case gateway.ErrSelfOrderNotAllowed.Error():
			return gateway.ErrSelfOrderNotAllowed
		case gateway.ErrOrderNotConfirmable.Error():
			return gateway.ErrOrderNotConfirmable
		case gateway.ErrOrderRequirementsIncomplete.Error():
			return gateway.ErrOrderRequirementsIncomplete
		case gateway.ErrOrderAlreadyPaymentPending.Error():
			return gateway.ErrOrderAlreadyPaymentPending
		case gateway.ErrOrderAlreadyFunded.Error():
			return gateway.ErrOrderAlreadyFunded
		case gateway.ErrOrderNotDeliverable.Error():
			return gateway.ErrOrderNotDeliverable
		case gateway.ErrOrderNotAcceptable.Error():
			return gateway.ErrOrderNotAcceptable
		case gateway.ErrOrderNotRevisionable.Error():
			return gateway.ErrOrderNotRevisionable
		case gateway.ErrOrderNotDisputable.Error():
			return gateway.ErrOrderNotDisputable
		case gateway.ErrOrderReleaseFailed.Error():
			return gateway.ErrOrderReleaseFailed
		case gateway.ErrConnectOnboardingIncomplete.Error():
			return gateway.ErrConnectOnboardingIncomplete
		default:
			return gateway.ErrInvalidGigState
		}
	case codes.AlreadyExists:
		return gateway.ErrFailedToCreateOrder
	default:
		return err
	}
}
