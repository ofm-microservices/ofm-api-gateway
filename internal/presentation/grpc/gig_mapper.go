package grpc

import (
	gateway "api-gateway/internal/domain"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	gigv1 "github.com/ofm-microseervices/ofm-common/proto/gig/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gigMapper struct {
	log logging.Logger
}

func newGigMapper(log logging.Logger) GigMapper {
	return &gigMapper{log: log}
}

func (m *gigMapper) ToCreateDraftRequest(req gateway.CreateGigDraftRequest) *gigv1.CreateDraftRequest {
	return &gigv1.CreateDraftRequest{FreelancerId: req.FreelancerID}
}

func (m *gigMapper) ToCreateDraftResponse(res *gigv1.CreateDraftResponse) *gateway.Gig {
	return m.toGig(res.GetGig())
}

func (m *gigMapper) ToUpdateBasicInfoRequest(req gateway.UpdateGigBasicInfoRequest) *gigv1.UpdateBasicInfoRequest {
	return &gigv1.UpdateBasicInfoRequest{
		GigId:        req.GigID,
		FreelancerId: req.FreelancerID,
		Title:        req.Title,
		Description:  req.Description,
		CategoryId:   req.CategoryID,
		Currency:     req.Currency,
	}
}

func (m *gigMapper) ToUpdateBasicInfoResponse(res *gigv1.UpdateBasicInfoResponse) *gateway.Gig {
	return m.toGig(res.GetGig())
}

func (m *gigMapper) ToReplacePackagesRequest(req gateway.ReplaceGigPackagesRequest) *gigv1.ReplacePackagesRequest {
	packages := make([]*gigv1.GigPackage, 0, len(req.Packages))
	for _, pkg := range req.Packages {
		packages = append(packages, &gigv1.GigPackage{
			Id:           pkg.ID,
			GigId:        pkg.GigID,
			Tier:         pkg.Tier,
			Description:  pkg.Description,
			DeliveryDays: pkg.DeliveryDays,
			PriceCents:   pkg.PriceCents,
			SortOrder:    pkg.SortOrder,
		})
	}

	return &gigv1.ReplacePackagesRequest{
		GigId:        req.GigID,
		FreelancerId: req.FreelancerID,
		Packages:     packages,
	}
}

func (m *gigMapper) ToReplacePackagesResponse(res *gigv1.ReplacePackagesResponse) *gateway.Gig {
	return m.toGig(res.GetGig())
}

func (m *gigMapper) ToReplaceQuestionsRequest(req gateway.ReplaceGigQuestionsRequest) *gigv1.ReplaceQuestionsRequest {
	questions := make([]*gigv1.GigQuestion, 0, len(req.Questions))
	for _, q := range req.Questions {
		questions = append(questions, &gigv1.GigQuestion{
			Id:        q.ID,
			GigId:     q.GigID,
			Content:   q.Content,
			SortOrder: q.SortOrder,
		})
	}

	return &gigv1.ReplaceQuestionsRequest{
		GigId:        req.GigID,
		FreelancerId: req.FreelancerID,
		Questions:    questions,
	}
}

func (m *gigMapper) ToReplaceQuestionsResponse(res *gigv1.ReplaceQuestionsResponse) *gateway.Gig {
	return m.toGig(res.GetGig())
}

func (m *gigMapper) ToReplaceMediaRequest(req gateway.ReplaceGigMediaRequest) *gigv1.ReplaceMediaRequest {
	files := make([]*gigv1.MediaUpload, 0, len(req.Files))
	for _, item := range req.Files {
		files = append(files, &gigv1.MediaUpload{
			Filename:    item.Filename,
			ContentType: item.ContentType,
			Data:        item.Data,
		})
	}

	return &gigv1.ReplaceMediaRequest{
		GigId:        req.GigID,
		FreelancerId: req.FreelancerID,
		Files:        files,
	}
}

func (m *gigMapper) ToReplaceMediaResponse(res *gigv1.ReplaceMediaResponse) *gateway.Gig {
	return m.toGig(res.GetGig())
}

func (m *gigMapper) ToGetDraftRequest(req gateway.GetGigDraftRequest) *gigv1.GetDraftRequest {
	return &gigv1.GetDraftRequest{GigId: req.GigID, FreelancerId: req.FreelancerID}
}

func (m *gigMapper) ToGetDraftResponse(res *gigv1.GetDraftResponse) *gateway.Gig {
	return m.toGig(res.GetGig())
}

func (m *gigMapper) ToPublishRequest(req gateway.PublishGigRequest) *gigv1.PublishRequest {
	return &gigv1.PublishRequest{GigId: req.GigID, FreelancerId: req.FreelancerID}
}

func (m *gigMapper) ToPublishResponse(res *gigv1.PublishResponse) *gateway.Gig {
	return m.toGig(res.GetGig())
}

func (m *gigMapper) ToError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return gateway.ErrFailedToCreateGig
	}

	switch st.Message() {
	case gateway.ErrInvalidGigID.Error():
		return gateway.ErrInvalidGigID
	case gateway.ErrInvalidFreelancerID.Error():
		return gateway.ErrInvalidFreelancerID
	case gateway.ErrInvalidTitle.Error():
		return gateway.ErrInvalidTitle
	case gateway.ErrInvalidDescription.Error():
		return gateway.ErrInvalidDescription
	case gateway.ErrInvalidCategoryID.Error():
		return gateway.ErrInvalidCategoryID
	case gateway.ErrInvalidCurrency.Error():
		return gateway.ErrInvalidCurrency
	case gateway.ErrInvalidPackageTier.Error():
		return gateway.ErrInvalidPackageTier
	case gateway.ErrInvalidPackageDescription.Error():
		return gateway.ErrInvalidPackageDescription
	case gateway.ErrInvalidPackageDeliveryDays.Error():
		return gateway.ErrInvalidPackageDeliveryDays
	case gateway.ErrInvalidPackagePriceCents.Error():
		return gateway.ErrInvalidPackagePriceCents
	case gateway.ErrInvalidQuestionContent.Error():
		return gateway.ErrInvalidQuestionContent
	case gateway.ErrInvalidMediaUpload.Error():
		return gateway.ErrInvalidMediaUpload
	case gateway.ErrInvalidPackageCount.Error():
		return gateway.ErrInvalidPackageCount
	case gateway.ErrGigNotFound.Error():
		return gateway.ErrGigNotFound
	case gateway.ErrGigDraftIncomplete.Error():
		return gateway.ErrGigDraftIncomplete
	case gateway.ErrGigAlreadyPublished.Error():
		return gateway.ErrGigAlreadyPublished
	case gateway.ErrInvalidGigState.Error():
		return gateway.ErrInvalidGigState
	default:
		if st.Code() == codes.NotFound {
			return gateway.ErrGigNotFound
		}
		if st.Code() == codes.InvalidArgument {
			return gateway.ErrFailedToUpdateGig
		}
		return gateway.ErrFailedToCreateGig
	}
}

func (m *gigMapper) toGig(res *gigv1.Gig) *gateway.Gig {
	if res == nil {
		return nil
	}

	gig := &gateway.Gig{
		GigID:                 res.GetGigId(),
		FreelancerID:          res.GetFreelancerId(),
		Title:                 res.GetTitle(),
		Description:           res.GetDescription(),
		CategoryID:            res.GetCategoryId(),
		Currency:              res.GetCurrency(),
		Status:                res.GetStatus(),
		BasicInfoCompleted:    res.GetBasicInfoCompleted(),
		PackagesCompleted:     res.GetPackagesCompleted(),
		RequirementsCompleted: res.GetRequirementsCompleted(),
		MediaCompleted:        res.GetMediaCompleted(),
		PictureFileID:         res.GetPictureFileId(),
		PublishedAt:           res.GetPublishedAt(),
		CreatedAt:             res.GetCreatedAt(),
		UpdatedAt:             res.GetUpdatedAt(),
	}

	if len(res.GetPackages()) > 0 {
		gig.Packages = make([]gateway.GigPackage, 0, len(res.GetPackages()))
		for _, pkg := range res.GetPackages() {
			gig.Packages = append(gig.Packages, gateway.GigPackage{
				ID:           pkg.GetId(),
				GigID:        pkg.GetGigId(),
				Tier:         pkg.GetTier(),
				Description:  pkg.GetDescription(),
				DeliveryDays: pkg.GetDeliveryDays(),
				PriceCents:   pkg.GetPriceCents(),
				SortOrder:    pkg.GetSortOrder(),
			})
		}
	}

	if len(res.GetQuestions()) > 0 {
		gig.Questions = make([]gateway.GigQuestion, 0, len(res.GetQuestions()))
		for _, q := range res.GetQuestions() {
			gig.Questions = append(gig.Questions, gateway.GigQuestion{
				ID:        q.GetId(),
				GigID:     q.GetGigId(),
				Content:   q.GetContent(),
				SortOrder: q.GetSortOrder(),
			})
		}
	}

	if len(res.GetMedia()) > 0 {
		gig.Media = make([]gateway.GigMedia, 0, len(res.GetMedia()))
		for _, item := range res.GetMedia() {
			gig.Media = append(gig.Media, gateway.GigMedia{
				ID:        item.GetFileId(),
				GigID:     item.GetGigId(),
				FileID:    item.GetFileId(),
				SortOrder: item.GetSortOrder(),
			})
		}
	}

	return gig
}
