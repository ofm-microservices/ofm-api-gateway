package grpc

import (
	"strings"

	gateway "api-gateway/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonv1 "github.com/ofm-microservices/ofm-common/proto/common/v1"
	orderwritev1 "github.com/ofm-microservices/ofm-common/proto/orderwrite/v1"
)

type orderPreviewMapper struct{}

func newOrderPreviewMapper() OrderPreviewMapper { return &orderPreviewMapper{} }

func (m *orderPreviewMapper) ToGetOrderPreviewByIDRequest(req gateway.GetOrderPreviewByIDRequest) *orderwritev1.GetOrderPreviewByIDRequest {
	role := commonv1.ParticipantRole_PARTICIPANT_ROLE_UNSPECIFIED
	switch req.Role {
	case gateway.ParticipantRoleCustomer:
		role = commonv1.ParticipantRole_PARTICIPANT_ROLE_BUYER
	case gateway.ParticipantRoleFreelancer:
		role = commonv1.ParticipantRole_PARTICIPANT_ROLE_SELLER
	}
	return &orderwritev1.GetOrderPreviewByIDRequest{
		OrderId: strings.TrimSpace(req.OrderID),
		UserId:  strings.TrimSpace(req.UserID),
		Role:    role,
	}
}

func (m *orderPreviewMapper) ToGetOrderPreviewByIDResponse(res *orderwritev1.GetOrderPreviewByIDResponse) *gateway.GetOrderPreviewByIDResult {
	if res == nil || res.GetOrder() == nil {
		return &gateway.GetOrderPreviewByIDResult{}
	}
	out := &gateway.GetOrderPreviewByIDResult{
		Order: &gateway.OrderPreview{
			OrderID:   res.GetOrder().GetOrderId(),
			CreatedAt: res.GetOrder().GetCreatedAt(),
			Status:    res.GetOrder().GetStatus(),
		},
	}
	if gig := res.GetGig(); gig != nil {
		out.Gig = &gateway.OrderPreviewGig{
			GigID:               gig.GetGigId(),
			Title:               gig.GetTitle(),
			PictureURL:          gig.GetPictureUrl(),
			PackageID:           gig.GetPackageId(),
			PackageTitle:        gig.GetPackageTitle(),
			PriceCents:          gig.GetPriceCents(),
			Currency:            gig.GetCurrency(),
			Description:         gig.GetDescription(),
			PackageDeliveryDays: gig.GetPackageDeliveryDays(),
		}
	}
	if customer := res.GetCustomer(); customer != nil {
		out.Customer = &gateway.OrderPreviewUser{
			UserID:      customer.GetUserId(),
			Username:    customer.GetUsername(),
			DisplayName: customer.GetDisplayName(),
			AvatarURL:   customer.GetAvatarUrl(),
		}
	}
	if freelancer := res.GetFreelancer(); freelancer != nil {
		out.Freelancer = &gateway.OrderPreviewUser{
			UserID:      freelancer.GetUserId(),
			Username:    freelancer.GetUsername(),
			DisplayName: freelancer.GetDisplayName(),
			AvatarURL:   freelancer.GetAvatarUrl(),
		}
	}
	return out
}

func (m *orderPreviewMapper) ToError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.InvalidArgument:
		if strings.Contains(strings.ToLower(st.Message()), "participant role") {
			return gateway.ErrInvalidParticipantRole
		}
		return gateway.ErrInvalidOrderID
	case codes.NotFound:
		return gateway.ErrOrderNotOwned
	default:
		return err
	}
}
