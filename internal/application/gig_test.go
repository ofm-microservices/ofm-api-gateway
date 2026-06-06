package service

import (
	"context"
	"errors"

	gateway "api-gateway/internal/domain"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
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
	slugReq      gateway.GetGigBySlugRequest
	previewReq   gateway.GetPreviewGigsByFreelancerUsernameRequest
	myGigsReq    gateway.GetMyGigsRequest
	publishReq   gateway.PublishGigRequest

	createRes    *gateway.Gig
	updateRes    *gateway.Gig
	packagesRes  *gateway.Gig
	questionsRes *gateway.Gig
	mediaRes     *gateway.Gig
	draftRes     *gateway.Gig
	slugRes      *gateway.Gig
	previewRes   *gateway.GigPreviewList
	myGigsRes    *gateway.GigPreviewPage
	publishRes   *gateway.Gig

	createErr    error
	updateErr    error
	packagesErr  error
	questionsErr error
	mediaErr     error
	draftErr     error
	slugErr      error
	previewErr   error
	myGigsErr    error
	publishErr   error
}

type reviewClientStub struct {
	reviewsReq       gateway.GetGigReviewsRequest
	summaryReq       gateway.GetGigReviewsSummaryRequest
	sellerSummaryReq gateway.GetUserRatingSummaryByUsernameRequest
	sellerReviewsReq gateway.ListSellerReviewsRequest
	reviewsRes       *gateway.GetGigReviewsResult
	summaryRes       *gateway.ReviewSummary
	sellerReviewsRes *gateway.ListSellerReviewsResult
	reviewsErr       error
	summaryErr       error
	sellerReviewsErr error
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

func (s *gigPublisherStub) GetBySlug(_ context.Context, req gateway.GetGigBySlugRequest) (*gateway.Gig, error) {
	s.slugReq = req
	return s.slugRes, s.slugErr
}

func (s *gigPublisherStub) GetPreviewGigsByFreelancerUsername(_ context.Context, req gateway.GetPreviewGigsByFreelancerUsernameRequest) (*gateway.GigPreviewList, error) {
	s.previewReq = req
	return s.previewRes, s.previewErr
}

func (s *gigPublisherStub) GetMyGigs(_ context.Context, req gateway.GetMyGigsRequest) (*gateway.GigPreviewPage, error) {
	s.myGigsReq = req
	return s.myGigsRes, s.myGigsErr
}

func (s *gigPublisherStub) Publish(_ context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error) {
	s.publishReq = req
	return s.publishRes, s.publishErr
}

type userProfileServiceStub struct {
	req gateway.GetUserProfileRequest
	res *gateway.UserProfile
	err error
}

func (s *userProfileServiceStub) GetUserProfile(_ context.Context, req gateway.GetUserProfileRequest) (*gateway.UserProfile, error) {
	s.req = req
	return s.res, s.err
}

func (s *reviewClientStub) CreateReview(context.Context, gateway.CreateReviewRequest) (*gateway.CreateReviewResult, error) {
	return nil, nil
}

func (s *reviewClientStub) GetGigReviews(_ context.Context, req gateway.GetGigReviewsRequest) (*gateway.GetGigReviewsResult, error) {
	s.reviewsReq = req
	return s.reviewsRes, s.reviewsErr
}

func (s *reviewClientStub) ListSellerReviews(_ context.Context, req gateway.ListSellerReviewsRequest) (*gateway.ListSellerReviewsResult, error) {
	s.sellerReviewsReq = req
	if s.sellerReviewsRes != nil {
		return s.sellerReviewsRes, s.sellerReviewsErr
	}
	return &gateway.ListSellerReviewsResult{Reviews: &gateway.ReviewList{}}, s.sellerReviewsErr
}

func (s *reviewClientStub) GetReviewsBySellerUsername(_ context.Context, req gateway.GetReviewsBySellerUsernameRequest) (*gateway.ListSellerReviewsResult, error) {
	s.sellerReviewsReq = gateway.ListSellerReviewsRequest{Username: req.Username, Cursor: req.Cursor}
	if s.sellerReviewsRes != nil {
		return s.sellerReviewsRes, s.sellerReviewsErr
	}
	return &gateway.ListSellerReviewsResult{Reviews: &gateway.ReviewList{}}, s.sellerReviewsErr
}

func (s *reviewClientStub) GetGigReviewsSummary(_ context.Context, req gateway.GetGigReviewsSummaryRequest) (*gateway.ReviewSummary, error) {
	s.summaryReq = req
	return s.summaryRes, s.summaryErr
}

func (s *reviewClientStub) GetUserRatingSummaryByUsername(_ context.Context, req gateway.GetUserRatingSummaryByUsernameRequest) (*gateway.ReviewSummary, error) {
	s.sellerSummaryReq = req
	return s.summaryRes, s.summaryErr
}

func (s *reviewClientStub) Close() error { return nil }

type userClientStub struct {
	userIDReq   string
	usernameReq string
	userRes     *gateway.User
	err         error
}

func (s *userClientStub) GetUserPreviewByID(_ context.Context, userID string) (*gateway.User, error) {
	s.userIDReq = userID
	return s.userRes, s.err
}

func (s *userClientStub) GetDetailedUserByUsername(_ context.Context, username string) (*gateway.User, error) {
	s.usernameReq = username
	return s.userRes, s.err
}

func (s *userClientStub) Close() error { return nil }

var _ = Describe("GigService", func() {
	var (
		pub    *gigPublisherStub
		review *reviewClientStub
		user   *userClientStub
		lg     logging.Logger
	)

	BeforeEach(func() {
		pub = &gigPublisherStub{}
		review = &reviewClientStub{}
		user = &userClientStub{}
		var err error
		lg, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	It("validates constructor dependencies", func() {
		svc, err := NewGig(nil, review, user, lg)
		Expect(svc).To(BeNil())
		Expect(err).To(MatchError(ErrNilGigClient))

		svc, err = NewGig(pub, nil, user, lg)
		Expect(svc).To(BeNil())
		Expect(err).To(MatchError(ErrNilReviewClient))

		svc, err = NewGig(pub, review, nil, lg)
		Expect(svc).To(BeNil())
		Expect(err).To(MatchError(ErrNilUserClient))

		svc, err = NewGig(pub, review, user, nil)
		Expect(svc).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("creates drafts with trimmed freelancer ids", func() {
		svc, err := NewGig(pub, review, user, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.createRes = &gateway.Gig{GigID: "gig-1"}

		res, err := svc.CreateDraft(context.Background(), gateway.CreateGigDraftRequest{FreelancerID: " freelancer-1 "})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.GigID).To(Equal("gig-1"))
		Expect(pub.createReq.FreelancerID).To(Equal("freelancer-1"))
	})

	It("rejects invalid basic info fields", func() {
		svc, err := NewGig(pub, review, user, lg)
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
		svc, err := NewGig(pub, review, user, lg)
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
		svc, err := NewGig(pub, review, user, lg)
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
		svc, err := NewGig(pub, review, user, lg)
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

		_, err = svc.Publish(context.Background(), gateway.PublishGigRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Username: "alex"})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.publishReq.GigID).To(Equal("gig-1"))
		Expect(pub.publishReq.Username).To(Equal("alex"))
	})

	It("assembles the public user profile from user, gig, and review services", func() {
		gigs := &gigPublisherStub{}
		reviews := &reviewClientStub{}
		users := &userClientStub{}
		users.userRes = &gateway.User{UserID: "user-1", Username: "alex", DisplayName: "Alex Tester"}
		gigs.previewRes = &gateway.GigPreviewList{
			Items:      []gateway.GigPreview{{GigID: "gig-1"}},
			Page:       2,
			Limit:      10,
			TotalPages: 5,
		}
		reviews.sellerReviewsRes = &gateway.ListSellerReviewsResult{
			Reviews: &gateway.ReviewList{
				Items:   []gateway.Review{{ReviewID: "review-1"}},
				Cursor:  "reviews-cursor",
				HasMore: false,
			},
		}

		svc, err := NewUserProfile(gigs, reviews, users, lg)
		Expect(err).NotTo(HaveOccurred())

		res, err := svc.GetUserProfile(context.Background(), gateway.GetUserProfileRequest{
			Username:      " alex ",
			GigsCursor:    " g ",
			ReviewsCursor: " r ",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.User.Username).To(Equal("alex"))
		Expect(gigs.previewReq.Username).To(Equal("alex"))
		Expect(gigs.previewReq.Page).To(Equal(int32(1)))
		Expect(gigs.previewReq.Limit).To(Equal(int32(10)))
		Expect(reviews.sellerReviewsReq.Username).To(Equal("alex"))
		Expect(reviews.sellerReviewsReq.Cursor).To(Equal("r"))
		Expect(res.Gigs.Items).To(HaveLen(1))
		Expect(res.Reviews.Items).To(HaveLen(1))
	})

	It("returns delegate failures unchanged", func() {
		svc, err := NewGig(pub, review, user, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.createErr = errors.New("boom")

		_, err = svc.CreateDraft(context.Background(), gateway.CreateGigDraftRequest{FreelancerID: "freelancer-1"})
		Expect(err).To(MatchError("boom"))
	})

	It("validates all gig base request fields", func() {
		svc, err := NewGig(pub, review, user, lg)
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
		svc, err := NewGig(pub, review, user, lg)
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
		svc, err := NewGig(pub, review, user, lg)
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
		svc, err := NewGig(pub, review, user, lg)
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

		_, err = svc.Publish(context.Background(), gateway.PublishGigRequest{GigID: " gig-1 ", FreelancerID: " freelancer-1 ", Username: " alex "})
		Expect(err).NotTo(HaveOccurred())
		Expect(pub.publishReq).To(Equal(gateway.PublishGigRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Username: "alex"}))
	})

	It("loads public gig pages with reviews summary and freelancer", func() {
		svc, err := NewGig(pub, review, user, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.slugRes = &gateway.Gig{GigID: "019e706c-616e-7473-9c1a-838c33b75013", Slug: "my-gig-019e706c-616e-7473-9c1a-838c33b75013"}
		review.reviewsRes = &gateway.GetGigReviewsResult{Reviews: &gateway.ReviewList{Items: []gateway.Review{{ReviewID: "review-1"}}}}
		review.summaryRes = &gateway.ReviewSummary{TotalReviews: 1}
		user.userRes = &gateway.User{UserID: "user-1", Username: "alex", DisplayName: "Alex Tester"}

		res, err := svc.GetBySlug(context.Background(), gateway.GetGigBySlugRequest{Username: "alex", Slug: "my-gig-019e706c-616e-7473-9c1a-838c33b75013"})
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Freelancer).NotTo(BeNil())
		Expect(res.Freelancer.Username).To(Equal("alex"))
		Expect(res.Freelancer.ReviewsSummary).NotTo(BeNil())
		Expect(res.Freelancer.ReviewsSummary.TotalReviews).To(Equal(int64(1)))
		Expect(review.reviewsReq.GigID).To(Equal("019e706c-616e-7473-9c1a-838c33b75013"))
		Expect(review.reviewsReq.Cursor).To(Equal(""))
		Expect(review.summaryReq.GigID).To(Equal("019e706c-616e-7473-9c1a-838c33b75013"))
		Expect(review.sellerSummaryReq.Username).To(Equal("alex"))
		Expect(user.usernameReq).To(Equal("alex"))
	})

	It("passes the page parameters through to review lookup", func() {
		svc, err := NewGig(pub, review, user, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.slugRes = &gateway.Gig{GigID: "019e706c-616e-7473-9c1a-838c33b75013", Slug: "my-gig-019e706c-616e-7473-9c1a-838c33b75013"}

		_, err = svc.GetBySlug(context.Background(), gateway.GetGigBySlugRequest{
			Username: "alex",
			Slug:     "my-gig-019e706c-616e-7473-9c1a-838c33b75013",
			Cursor:   "abc",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(review.reviewsReq.Cursor).To(Equal("abc"))
	})

	It("returns not found when the username does not match the gig owner", func() {
		svc, err := NewGig(pub, review, user, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.slugRes = &gateway.Gig{
			GigID:        "019e706c-616e-7473-9c1a-838c33b75013",
			FreelancerID: "freelancer-1",
			Slug:         "my-gig-019e706c-616e-7473-9c1a-838c33b75013",
		}
		user.userRes = &gateway.User{UserID: "freelancer-2", Username: "alex"}

		res, err := svc.GetBySlug(context.Background(), gateway.GetGigBySlugRequest{
			Username: "alex",
			Slug:     "my-gig-019e706c-616e-7473-9c1a-838c33b75013",
		})
		Expect(res).To(BeNil())
		Expect(err).To(MatchError(gateway.ErrGigNotFound))
	})

	It("returns not found when the requested slug is not canonical", func() {
		svc, err := NewGig(pub, review, user, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.slugRes = &gateway.Gig{
			GigID:        "019e706c-616e-7473-9c1a-838c33b75013",
			FreelancerID: "019e4767-c408-74d7-8c4a-3b7fcb2ab3cf",
			Slug:         "my-gig-019e706c-616e-7473-9c1a-838c33b75013",
		}
		user.userRes = &gateway.User{UserID: "019e4767-c408-74d7-8c4a-3b7fcb2ab3cf", Username: "alex"}

		res, err := svc.GetBySlug(context.Background(), gateway.GetGigBySlugRequest{
			Username: "alex",
			Slug:     "my-kgig-019e706c-616e-7473-9c1a-838c33b75013",
		})
		Expect(res).To(BeNil())
		Expect(err).To(MatchError(gateway.ErrGigNotFound))
	})

	It("returns the gig when user and review enrichment fails", func() {
		svc, err := NewGig(pub, review, user, lg)
		Expect(err).NotTo(HaveOccurred())
		pub.slugRes = &gateway.Gig{
			GigID:        "019e706c-616e-7473-9c1a-838c33b75013",
			FreelancerID: "019e4767-c408-74d7-8c4a-3b7fcb2ab3cf",
			Slug:         "my-gig-019e706c-616e-7473-9c1a-838c33b75013",
			Media:        []gateway.GigMedia{{GigID: "019e706c-616e-7473-9c1a-838c33b75013", URL: "https://example.com/media.jpg"}},
		}
		user.err = errors.New("user down")
		review.reviewsErr = errors.New("reviews down")
		review.summaryErr = errors.New("summary down")

		res, err := svc.GetBySlug(context.Background(), gateway.GetGigBySlugRequest{
			Username: "alex",
			Slug:     "my-gig-019e706c-616e-7473-9c1a-838c33b75013",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.Freelancer).To(BeNil())
		Expect(res.Reviews).To(BeNil())
		Expect(res.ReviewsSummary).To(BeNil())
		Expect(res.Media).To(HaveLen(1))
		Expect(res.Media[0].URL).To(Equal("https://example.com/media.jpg"))
	})
})
