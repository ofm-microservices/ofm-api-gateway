package service

import (
	gateway "api-gateway/internal/domain"
	"context"
	"strings"

	"github.com/ofm-microseervices/ofm-common/pkg/logging"
)

type gigService struct {
	client GigPublisher
	log    Logger
}

// NewGig constructs the application service responsible for managing gig
// drafts through the gig-service boundary.
func NewGig(client GigPublisher, log Logger) (GigService, error) {
	if client == nil {
		return nil, ErrNilGigClient
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &gigService{
		client: client,
		log:    log.With(logging.String("module", "gig-application")),
	}, nil
}

func (s *gigService) CreateDraft(ctx context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error) {
	freelancerID := strings.TrimSpace(req.FreelancerID)
	if freelancerID == "" {
		return nil, gateway.ErrInvalidFreelancerID
	}

	return s.client.CreateDraft(ctx, gateway.CreateGigDraftRequest{FreelancerID: freelancerID})
}

func (s *gigService) UpdateBasicInfo(ctx context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error) {
	normalized, err := normalizeGigBaseRequest(req.GigID, req.FreelancerID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Title) == "" {
		return nil, gateway.ErrInvalidTitle
	}
	if strings.TrimSpace(req.Description) == "" {
		return nil, gateway.ErrInvalidDescription
	}
	if req.CategoryID <= 0 {
		return nil, gateway.ErrInvalidCategoryID
	}
	if strings.TrimSpace(req.Currency) == "" {
		return nil, gateway.ErrInvalidCurrency
	}

	return s.client.UpdateBasicInfo(ctx, gateway.UpdateGigBasicInfoRequest{
		GigID:        normalized.gigID,
		FreelancerID: normalized.freelancerID,
		Title:        strings.TrimSpace(req.Title),
		Description:  strings.TrimSpace(req.Description),
		CategoryID:   req.CategoryID,
		Currency:     strings.TrimSpace(req.Currency),
	})
}

func (s *gigService) ReplacePackages(ctx context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error) {
	normalized, err := normalizeGigBaseRequest(req.GigID, req.FreelancerID)
	if err != nil {
		return nil, err
	}
	packages, err := normalizeGigPackages(req.Packages)
	if err != nil {
		return nil, err
	}

	return s.client.ReplacePackages(ctx, gateway.ReplaceGigPackagesRequest{
		GigID:        normalized.gigID,
		FreelancerID: normalized.freelancerID,
		Packages:     packages,
	})
}

func (s *gigService) ReplaceQuestions(ctx context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error) {
	normalized, err := normalizeGigBaseRequest(req.GigID, req.FreelancerID)
	if err != nil {
		return nil, err
	}
	questions := make([]gateway.GigQuestion, 0, len(req.Questions))
	for _, q := range req.Questions {
		if strings.TrimSpace(q.Content) == "" {
			return nil, gateway.ErrInvalidQuestionContent
		}
		questions = append(questions, gateway.GigQuestion{
			ID:        strings.TrimSpace(q.ID),
			GigID:     normalized.gigID,
			Content:   strings.TrimSpace(q.Content),
			SortOrder: q.SortOrder,
		})
	}

	return s.client.ReplaceQuestions(ctx, gateway.ReplaceGigQuestionsRequest{
		GigID:        normalized.gigID,
		FreelancerID: normalized.freelancerID,
		Questions:    questions,
	})
}

func (s *gigService) ReplaceMedia(ctx context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error) {
	normalized, err := normalizeGigBaseRequest(req.GigID, req.FreelancerID)
	if err != nil {
		return nil, err
	}
	files := make([]gateway.GigMediaUpload, 0, len(req.Files))
	for _, item := range req.Files {
		if strings.TrimSpace(item.Filename) == "" || strings.TrimSpace(item.ContentType) == "" || len(item.Data) == 0 {
			return nil, gateway.ErrInvalidMediaUpload
		}
		files = append(files, gateway.GigMediaUpload{
			Filename:    strings.TrimSpace(item.Filename),
			ContentType: strings.TrimSpace(item.ContentType),
			Data:        item.Data,
		})
	}

	return s.client.ReplaceMedia(ctx, gateway.ReplaceGigMediaRequest{
		GigID:        normalized.gigID,
		FreelancerID: normalized.freelancerID,
		Files:        files,
	})
}

func (s *gigService) GetDraft(ctx context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error) {
	normalized, err := normalizeGigBaseRequest(req.GigID, req.FreelancerID)
	if err != nil {
		return nil, err
	}

	return s.client.GetDraft(ctx, gateway.GetGigDraftRequest{
		GigID:        normalized.gigID,
		FreelancerID: normalized.freelancerID,
	})
}

func (s *gigService) Publish(ctx context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error) {
	normalized, err := normalizeGigBaseRequest(req.GigID, req.FreelancerID)
	if err != nil {
		return nil, err
	}

	return s.client.Publish(ctx, gateway.PublishGigRequest{
		GigID:        normalized.gigID,
		FreelancerID: normalized.freelancerID,
	})
}

type normalizedGigRequest struct {
	gigID        string
	freelancerID string
}

func normalizeGigBaseRequest(gigID, freelancerID string) (normalizedGigRequest, error) {
	gigID = strings.TrimSpace(gigID)
	if gigID == "" {
		return normalizedGigRequest{}, gateway.ErrInvalidGigID
	}
	freelancerID = strings.TrimSpace(freelancerID)
	if freelancerID == "" {
		return normalizedGigRequest{}, gateway.ErrInvalidFreelancerID
	}

	return normalizedGigRequest{gigID: gigID, freelancerID: freelancerID}, nil
}

func normalizeGigPackages(packages []gateway.GigPackage) ([]gateway.GigPackage, error) {
	if len(packages) == 0 || len(packages) > 3 {
		return nil, gateway.ErrInvalidPackageCount
	}

	allowed := map[string]struct{}{
		gateway.TierBasic:    {},
		gateway.TierStandard: {},
		gateway.TierPremium:  {},
	}
	expectedOrder := []string{gateway.TierBasic, gateway.TierStandard, gateway.TierPremium}
	out := make([]gateway.GigPackage, 0, len(packages))

	for i, pkg := range packages {
		tier := strings.TrimSpace(pkg.Tier)
		if _, ok := allowed[tier]; !ok {
			return nil, gateway.ErrInvalidPackageTier
		}
		if tier != expectedOrder[i] {
			return nil, gateway.ErrInvalidPackageTier
		}
		if strings.TrimSpace(pkg.Description) == "" {
			return nil, gateway.ErrInvalidPackageDescription
		}
		if pkg.DeliveryDays <= 0 {
			return nil, gateway.ErrInvalidPackageDeliveryDays
		}
		if pkg.PriceCents <= 0 {
			return nil, gateway.ErrInvalidPackagePriceCents
		}
		out = append(out, gateway.GigPackage{
			ID:           strings.TrimSpace(pkg.ID),
			GigID:        strings.TrimSpace(pkg.GigID),
			Tier:         tier,
			Description:  strings.TrimSpace(pkg.Description),
			DeliveryDays: pkg.DeliveryDays,
			PriceCents:   pkg.PriceCents,
			SortOrder:    pkg.SortOrder,
		})
	}

	return out, nil
}
