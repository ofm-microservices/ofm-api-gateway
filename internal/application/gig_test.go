package service

import (
	"context"
	"errors"

	gateway "api-gateway/internal/domain"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type gigPublisherStub struct {
	createReq    gateway.CreateGigDraftRequest
	updateReq    gateway.UpdateGigBasicInfoRequest
	packagesReq  gateway.ReplaceGigPackagesRequest
	questionsReq gateway.ReplaceGigQuestionsRequest
	mediaReq     gateway.ReplaceGigMediaRequest
	draftReq     gateway.GetGigDraftRequest
	publishReq   gateway.PublishGigRequest

	createRes    *gateway.Gig
	updateRes    *gateway.Gig
	packagesRes  *gateway.Gig
	questionsRes *gateway.Gig
	mediaRes     *gateway.Gig
	draftRes     *gateway.Gig
	publishRes   *gateway.Gig

	createErr    error
	updateErr    error
	packagesErr  error
	questionsErr error
	mediaErr     error
	draftErr     error
	publishErr   error
}

func (s *gigPublisherStub) CreateDraft(_ context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error) {
	s.createReq = req
	return s.createRes, s.createErr
}

func (s *gigPublisherStub) UpdateBasicInfo(_ context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error) {
	s.updateReq = req
	return s.updateRes, s.updateErr
}

func (s *gigPublisherStub) ReplacePackages(_ context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error) {
	s.packagesReq = req
	return s.packagesRes, s.packagesErr
}

func (s *gigPublisherStub) ReplaceQuestions(_ context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error) {
	s.questionsReq = req
	return s.questionsRes, s.questionsErr
}

func (s *gigPublisherStub) ReplaceMedia(_ context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error) {
	s.mediaReq = req
	return s.mediaRes, s.mediaErr
}

func (s *gigPublisherStub) GetDraft(_ context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error) {
	s.draftReq = req
	return s.draftRes, s.draftErr
}

func (s *gigPublisherStub) Publish(_ context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error) {
	s.publishReq = req
	return s.publishRes, s.publishErr
}

var _ = Describe("GigService", func() {
	var (
		pub *gigPublisherStub
		lg  logging.Logger
	)

	BeforeEach(func() {
		pub = &gigPublisherStub{}
		var err error
		lg, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	It("validates constructor dependencies", func() {
		svc, err := NewGig(nil, lg)
		Expect(svc).To(BeNil())
		Expect(err).To(MatchError(ErrNilGigClient))

		svc, err = NewGig(pub, nil)
		Expect(svc).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("creates drafts with trimmed freelancer ids", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.createRes = &gateway.Gig{GigID: "gig-1"}

		res, err := svc.CreateDraft(context.Background(), gateway.CreateGigDraftRequest{FreelancerID: " freelancer-1 "})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.GigID).To(Equal("gig-1"))
		Expect(pub.createReq.FreelancerID).To(Equal("freelancer-1"))
	})

	It("rejects invalid basic info fields", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())

		_, err = svc.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{GigID: "gig-1", FreelancerID: "freelancer-1"})
		Expect(err).To(MatchError(gateway.ErrInvalidTitle))

		_, err = svc.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Title:        "Title",
			Description:  "Description",
			CategoryID:   1,
		})
		Expect(err).To(MatchError(gateway.ErrInvalidCurrency))
	})

	It("normalizes and forwards valid basic info", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.updateRes = &gateway.Gig{GigID: "gig-1"}

		res, err := svc.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{
			GigID:        " gig-1 ",
			FreelancerID: " freelancer-1 ",
			Title:        " Title ",
			Description:  " Description ",
			CategoryID:   7,
			Currency:     " usd ",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.GigID).To(Equal("gig-1"))
		Expect(pub.updateReq).To(Equal(gateway.UpdateGigBasicInfoRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Title:        "Title",
			Description:  "Description",
			CategoryID:   7,
			Currency:     "usd",
		}))
	})

	It("validates package replacement rules", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Packages: []gateway.GigPackage{{
				Tier:         gateway.TierStandard,
				Description:  "Good",
				DeliveryDays: 2,
				PriceCents:   2000,
			}},
		})
		Expect(err).To(MatchError(gateway.ErrInvalidPackageTier))

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Packages: []gateway.GigPackage{
				{Tier: gateway.TierBasic, Description: "Basic", DeliveryDays: 1, PriceCents: 1000},
				{Tier: gateway.TierStandard, Description: "Standard", DeliveryDays: 2, PriceCents: 2000},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.packagesReq.Packages).To(HaveLen(2))
	})

	It("forwards question, media, draft, and publish operations", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.questionsRes = &gateway.Gig{GigID: "gig-1"}
		pub.mediaRes = &gateway.Gig{GigID: "gig-1"}
		pub.draftRes = &gateway.Gig{GigID: "gig-1"}
		pub.publishRes = &gateway.Gig{GigID: "gig-1"}

		_, err = svc.ReplaceQuestions(context.Background(), gateway.ReplaceGigQuestionsRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Questions:    []gateway.GigQuestion{{Content: "Question?"}},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.questionsReq.Questions).To(HaveLen(1))

		_, err = svc.ReplaceMedia(context.Background(), gateway.ReplaceGigMediaRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Files:        []gateway.GigMediaUpload{{Filename: "cover.png", ContentType: "image/png", Data: []byte("cover")}},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.mediaReq.Files).To(HaveLen(1))

		_, err = svc.GetDraft(context.Background(), gateway.GetGigDraftRequest{GigID: "gig-1", FreelancerID: "freelancer-1"})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.draftReq.FreelancerID).To(Equal("freelancer-1"))

		_, err = svc.Publish(context.Background(), gateway.PublishGigRequest{GigID: "gig-1", FreelancerID: "freelancer-1"})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.publishReq.GigID).To(Equal("gig-1"))
	})

	It("returns delegate failures unchanged", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.createErr = errors.New("boom")

		_, err = svc.CreateDraft(context.Background(), gateway.CreateGigDraftRequest{FreelancerID: "freelancer-1"})
		Expect(err).To(MatchError("boom"))
	})

	It("validates all gig base request fields", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())

		_, err = svc.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{
			GigID:        " ",
			FreelancerID: "freelancer-1",
			Title:        "Title",
			Description:  "Description",
			CategoryID:   1,
			Currency:     "usd",
		})
		Expect(err).To(MatchError(gateway.ErrInvalidGigID))

		_, err = svc.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{
			GigID:        "gig-1",
			FreelancerID: " ",
			Title:        "Title",
			Description:  "Description",
			CategoryID:   1,
			Currency:     "usd",
		})
		Expect(err).To(MatchError(gateway.ErrInvalidFreelancerID))
	})

	It("validates every package field", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Packages:     nil,
		})
		Expect(err).To(MatchError(gateway.ErrInvalidPackageCount))

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Packages: []gateway.GigPackage{
				{Tier: gateway.TierBasic, Description: "Basic", DeliveryDays: 1, PriceCents: 1000},
				{Tier: "wrong", Description: "Standard", DeliveryDays: 2, PriceCents: 2000},
			},
		})
		Expect(err).To(MatchError(gateway.ErrInvalidPackageTier))

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Packages: []gateway.GigPackage{
				{Tier: gateway.TierBasic, Description: "", DeliveryDays: 1, PriceCents: 1000},
				{Tier: gateway.TierStandard, Description: "Standard", DeliveryDays: 2, PriceCents: 2000},
			},
		})
		Expect(err).To(MatchError(gateway.ErrInvalidPackageDescription))

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Packages: []gateway.GigPackage{
				{Tier: gateway.TierBasic, Description: "Basic", DeliveryDays: 0, PriceCents: 1000},
				{Tier: gateway.TierStandard, Description: "Standard", DeliveryDays: 2, PriceCents: 2000},
			},
		})
		Expect(err).To(MatchError(gateway.ErrInvalidPackageDeliveryDays))

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Packages: []gateway.GigPackage{
				{Tier: gateway.TierBasic, Description: "Basic", DeliveryDays: 1, PriceCents: 0},
				{Tier: gateway.TierStandard, Description: "Standard", DeliveryDays: 2, PriceCents: 2000},
			},
		})
		Expect(err).To(MatchError(gateway.ErrInvalidPackagePriceCents))
	})

	It("validates question and media payloads", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())

		_, err = svc.ReplaceQuestions(context.Background(), gateway.ReplaceGigQuestionsRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Questions:    []gateway.GigQuestion{{Content: " "}},
		})
		Expect(err).To(MatchError(gateway.ErrInvalidQuestionContent))

		_, err = svc.ReplaceMedia(context.Background(), gateway.ReplaceGigMediaRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Files:        []gateway.GigMediaUpload{{Filename: "", ContentType: "image/png", Data: []byte("x")}},
		})
		Expect(err).To(MatchError(gateway.ErrInvalidMediaUpload))
	})

	It("forwards cleaned gig state for all operations", func() {
		svc, err := NewGig(pub, lg)
		Expect(err).NotTo(HaveOccurred())

		pub.updateRes = &gateway.Gig{GigID: "gig-1"}
		pub.packagesRes = &gateway.Gig{GigID: "gig-1"}
		pub.questionsRes = &gateway.Gig{GigID: "gig-1"}
		pub.mediaRes = &gateway.Gig{GigID: "gig-1"}
		pub.draftRes = &gateway.Gig{GigID: "gig-1"}
		pub.publishRes = &gateway.Gig{GigID: "gig-1"}

		_, err = svc.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{
			GigID:        " gig-1 ",
			FreelancerID: " freelancer-1 ",
			Title:        " Title ",
			Description:  " Description ",
			CategoryID:   7,
			Currency:     " usd ",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.updateReq).To(Equal(gateway.UpdateGigBasicInfoRequest{
			GigID:        "gig-1",
			FreelancerID: "freelancer-1",
			Title:        "Title",
			Description:  "Description",
			CategoryID:   7,
			Currency:     "usd",
		}))

		_, err = svc.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{
			GigID:        " gig-1 ",
			FreelancerID: " freelancer-1 ",
			Packages: []gateway.GigPackage{
				{Tier: gateway.TierBasic, Description: " Basic ", DeliveryDays: 1, PriceCents: 1000},
				{Tier: gateway.TierStandard, Description: " Standard ", DeliveryDays: 2, PriceCents: 2000},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.packagesReq.Packages[0].Description).To(Equal("Basic"))

		_, err = svc.ReplaceQuestions(context.Background(), gateway.ReplaceGigQuestionsRequest{
			GigID:        " gig-1 ",
			FreelancerID: " freelancer-1 ",
			Questions:    []gateway.GigQuestion{{ID: " q-1 ", GigID: " old ", Content: " Question? ", SortOrder: 1}},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.questionsReq.Questions[0].ID).To(Equal("q-1"))
		Expect(pub.questionsReq.Questions[0].GigID).To(Equal("gig-1"))

		_, err = svc.ReplaceMedia(context.Background(), gateway.ReplaceGigMediaRequest{
			GigID:        " gig-1 ",
			FreelancerID: " freelancer-1 ",
			Files:        []gateway.GigMediaUpload{{Filename: " cover.png ", ContentType: " image/png ", Data: []byte("cover")}},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.mediaReq.Files[0].Filename).To(Equal("cover.png"))

		_, err = svc.GetDraft(context.Background(), gateway.GetGigDraftRequest{GigID: " gig-1 ", FreelancerID: " freelancer-1 "})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.draftReq).To(Equal(gateway.GetGigDraftRequest{GigID: "gig-1", FreelancerID: "freelancer-1"}))

		_, err = svc.Publish(context.Background(), gateway.PublishGigRequest{GigID: " gig-1 ", FreelancerID: " freelancer-1 "})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.publishReq).To(Equal(gateway.PublishGigRequest{GigID: "gig-1", FreelancerID: "freelancer-1"}))
	})
})
