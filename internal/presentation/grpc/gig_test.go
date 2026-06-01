package grpc

import (
	"context"
	"errors"

	gateway "api-gateway/internal/domain"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	gigv1 "github.com/ofm-microservices/ofm-common/proto/gig/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeGigCommandServiceClient struct {
	createReq    *gigv1.CreateDraftRequest
	updateReq    *gigv1.UpdateBasicInfoRequest
	packagesReq  *gigv1.ReplacePackagesRequest
	questionsReq *gigv1.ReplaceQuestionsRequest
	mediaReq     *gigv1.ReplaceMediaRequest
	draftReq     *gigv1.GetDraftRequest
	slugReq      *gigv1.GetGigBySlugRequest
	publishReq   *gigv1.PublishRequest

	createRes    *gigv1.CreateDraftResponse
	updateRes    *gigv1.UpdateBasicInfoResponse
	packagesRes  *gigv1.ReplacePackagesResponse
	questionsRes *gigv1.ReplaceQuestionsResponse
	mediaRes     *gigv1.ReplaceMediaResponse
	draftRes     *gigv1.GetDraftResponse
	slugRes      *gigv1.GetGigBySlugResponse
	publishRes   *gigv1.PublishResponse

	err error
}

func (f *fakeGigCommandServiceClient) CreateDraft(_ context.Context, req *gigv1.CreateDraftRequest, _ ...grpc.CallOption) (*gigv1.CreateDraftResponse, error) {
	f.createReq = req
	return f.createRes, f.err
}

func (f *fakeGigCommandServiceClient) UpdateBasicInfo(_ context.Context, req *gigv1.UpdateBasicInfoRequest, _ ...grpc.CallOption) (*gigv1.UpdateBasicInfoResponse, error) {
	f.updateReq = req
	return f.updateRes, f.err
}

func (f *fakeGigCommandServiceClient) ReplacePackages(_ context.Context, req *gigv1.ReplacePackagesRequest, _ ...grpc.CallOption) (*gigv1.ReplacePackagesResponse, error) {
	f.packagesReq = req
	return f.packagesRes, f.err
}

func (f *fakeGigCommandServiceClient) ReplaceQuestions(_ context.Context, req *gigv1.ReplaceQuestionsRequest, _ ...grpc.CallOption) (*gigv1.ReplaceQuestionsResponse, error) {
	f.questionsReq = req
	return f.questionsRes, f.err
}

func (f *fakeGigCommandServiceClient) ReplaceMedia(_ context.Context, req *gigv1.ReplaceMediaRequest, _ ...grpc.CallOption) (*gigv1.ReplaceMediaResponse, error) {
	f.mediaReq = req
	return f.mediaRes, f.err
}

func (f *fakeGigCommandServiceClient) GetDraft(_ context.Context, req *gigv1.GetDraftRequest, _ ...grpc.CallOption) (*gigv1.GetDraftResponse, error) {
	f.draftReq = req
	return f.draftRes, f.err
}

func (f *fakeGigCommandServiceClient) GetGigBySlug(_ context.Context, req *gigv1.GetGigBySlugRequest, _ ...grpc.CallOption) (*gigv1.GetGigBySlugResponse, error) {
	f.slugReq = req
	return f.slugRes, f.err
}

func (f *fakeGigCommandServiceClient) Publish(_ context.Context, req *gigv1.PublishRequest, _ ...grpc.CallOption) (*gigv1.PublishResponse, error) {
	f.publishReq = req
	return f.publishRes, f.err
}

func (f *fakeGigCommandServiceClient) GetOrderStartSnapshot(_ context.Context, req *gigv1.GetOrderStartSnapshotRequest, _ ...grpc.CallOption) (*gigv1.GetOrderStartSnapshotResponse, error) {
	return &gigv1.GetOrderStartSnapshotResponse{}, f.err
}

var _ = Describe("GigMapper", func() {
	It("maps requests and responses", func() {
		mapr := newGigMapper(logging.Logger(nil))
		req := mapr.ToCreateDraftRequest(gateway.CreateGigDraftRequest{FreelancerID: "freelancer-1"})
		Expect(req.GetFreelancerId()).To(Equal("freelancer-1"))

		res := mapr.ToCreateDraftResponse(&gigv1.CreateDraftResponse{
			Gig: &gigv1.Gig{
				GigId:        "gig-1",
				FreelancerId: "freelancer-1",
				Status:       "draft",
			},
		})
		Expect(res.GigID).To(Equal("gig-1"))
	})

	It("maps every request and response shape", func() {
		mapr := newGigMapper(logging.Logger(nil))
		now := "2026-05-10T15:00:00Z"
		pkg := &gigv1.GigPackage{Id: "pkg-1", GigId: "gig-1", Tier: gateway.TierBasic, Description: "basic", DeliveryDays: 1, PriceCents: 100, SortOrder: 1}
		question := &gigv1.GigQuestion{Id: "q-1", GigId: "gig-1", Content: "question", SortOrder: 1}
		media := &gigv1.GigMedia{FileId: "file-1", GigId: "gig-1", SortOrder: 1}
		gig := &gigv1.Gig{
			GigId:                 "gig-1",
			FreelancerId:          "freelancer-1",
			Slug:                  "my-gig-gig-1",
			Title:                 "Title",
			Description:           "Description",
			CategoryId:            1,
			Currency:              "usd",
			Status:                "draft",
			BasicInfoCompleted:    true,
			PackagesCompleted:     true,
			RequirementsCompleted: true,
			MediaCompleted:        true,
			PictureFileId:         "file-1",
			PublishedAt:           now,
			CreatedAt:             now,
			UpdatedAt:             now,
			Packages:              []*gigv1.GigPackage{pkg},
			Questions:             []*gigv1.GigQuestion{question},
			Media:                 []*gigv1.GigMedia{media},
		}

		Expect(mapr.ToUpdateBasicInfoRequest(gateway.UpdateGigBasicInfoRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Title: "Title", Description: "Description", CategoryID: 1, Currency: "usd"}).GetGigId()).To(Equal("gig-1"))
		Expect(mapr.ToReplacePackagesRequest(gateway.ReplaceGigPackagesRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Packages: []gateway.GigPackage{{ID: "pkg-1", GigID: "gig-1", Tier: gateway.TierBasic, Description: "basic", DeliveryDays: 1, PriceCents: 100, SortOrder: 1}}}).GetPackages()).To(HaveLen(1))
		Expect(mapr.ToReplaceQuestionsRequest(gateway.ReplaceGigQuestionsRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Questions: []gateway.GigQuestion{{ID: "q-1", GigID: "gig-1", Content: "question", SortOrder: 1}}}).GetQuestions()).To(HaveLen(1))
		Expect(mapr.ToReplaceMediaRequest(gateway.ReplaceGigMediaRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Files: []gateway.GigMediaUpload{{Filename: "cover.png", ContentType: "image/png", Data: []byte("x")}}}).GetFiles()).To(HaveLen(1))
		Expect(mapr.ToGetDraftRequest(gateway.GetGigDraftRequest{GigID: "gig-1", FreelancerID: "freelancer-1"}).GetGigId()).To(Equal("gig-1"))
		Expect(mapr.ToGetBySlugRequest(gateway.GetGigBySlugRequest{Slug: "my-gig-gig-1"}).GetSlug()).To(Equal("my-gig-gig-1"))
		Expect(mapr.ToPublishRequest(gateway.PublishGigRequest{GigID: "gig-1", FreelancerID: "freelancer-1"}).GetGigId()).To(Equal("gig-1"))

		Expect(mapr.ToUpdateBasicInfoResponse(&gigv1.UpdateBasicInfoResponse{Gig: gig}).GigID).To(Equal("gig-1"))
		Expect(mapr.ToReplacePackagesResponse(&gigv1.ReplacePackagesResponse{Gig: gig}).Packages).To(HaveLen(1))
		Expect(mapr.ToReplaceQuestionsResponse(&gigv1.ReplaceQuestionsResponse{Gig: gig}).Questions).To(HaveLen(1))
		Expect(mapr.ToReplaceMediaResponse(&gigv1.ReplaceMediaResponse{Gig: gig}).Media).To(HaveLen(1))
		Expect(mapr.ToGetDraftResponse(&gigv1.GetDraftResponse{Gig: gig}).GigID).To(Equal("gig-1"))
		Expect(mapr.ToGetBySlugResponse(&gigv1.GetGigBySlugResponse{Gig: gig}).GigID).To(Equal("gig-1"))
		Expect(mapr.ToPublishResponse(&gigv1.PublishResponse{Gig: gig}).GigID).To(Equal("gig-1"))
		Expect(mapr.ToCreateDraftResponse(&gigv1.CreateDraftResponse{Gig: gig}).GigID).To(Equal("gig-1"))
		Expect(mapr.ToCreateDraftResponse(nil)).To(BeNil())
		Expect(mapr.ToUpdateBasicInfoResponse(nil)).To(BeNil())
		Expect(mapr.ToReplacePackagesResponse(nil)).To(BeNil())
		Expect(mapr.ToReplaceQuestionsResponse(nil)).To(BeNil())
		Expect(mapr.ToReplaceMediaResponse(nil)).To(BeNil())
		Expect(mapr.ToGetDraftResponse(nil)).To(BeNil())
		Expect(mapr.ToPublishResponse(nil)).To(BeNil())
	})

	It("maps errors", func() {
		mapr := newGigMapper(logging.Logger(nil))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidGigID.Error()))).To(MatchError(gateway.ErrInvalidGigID))
		Expect(mapr.ToError(status.Error(codes.NotFound, "missing"))).To(MatchError(gateway.ErrGigNotFound))
		Expect(mapr.ToError(errors.New("boom"))).To(MatchError(gateway.ErrFailedToCreateGig))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidFreelancerID.Error()))).To(MatchError(gateway.ErrInvalidFreelancerID))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidTitle.Error()))).To(MatchError(gateway.ErrInvalidTitle))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidDescription.Error()))).To(MatchError(gateway.ErrInvalidDescription))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidCategoryID.Error()))).To(MatchError(gateway.ErrInvalidCategoryID))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidCurrency.Error()))).To(MatchError(gateway.ErrInvalidCurrency))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidPackageTier.Error()))).To(MatchError(gateway.ErrInvalidPackageTier))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidPackageDescription.Error()))).To(MatchError(gateway.ErrInvalidPackageDescription))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidPackageDeliveryDays.Error()))).To(MatchError(gateway.ErrInvalidPackageDeliveryDays))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidPackagePriceCents.Error()))).To(MatchError(gateway.ErrInvalidPackagePriceCents))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidQuestionContent.Error()))).To(MatchError(gateway.ErrInvalidQuestionContent))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidMediaUpload.Error()))).To(MatchError(gateway.ErrInvalidMediaUpload))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidPackageCount.Error()))).To(MatchError(gateway.ErrInvalidPackageCount))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrGigDraftIncomplete.Error()))).To(MatchError(gateway.ErrGigDraftIncomplete))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrGigAlreadyPublished.Error()))).To(MatchError(gateway.ErrGigAlreadyPublished))
		Expect(mapr.ToError(status.Error(codes.FailedPrecondition, gateway.ErrConnectOnboardingIncomplete.Error()))).To(MatchError(gateway.ErrConnectOnboardingIncomplete))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, gateway.ErrInvalidGigState.Error()))).To(MatchError(gateway.ErrInvalidGigState))
		Expect(mapr.ToError(status.Error(codes.InvalidArgument, "other"))).To(MatchError(gateway.ErrFailedToUpdateGig))
	})
})

var _ = Describe("GigClient", func() {
	It("validates constructor input", func() {
		cl, err := NewGigClient(GigServiceConfig{}, logging.Logger(nil))
		Expect(cl).To(BeNil())
		Expect(err).To(MatchError(ErrEmptyGigServiceAddress))

		cl, err = NewGigClient(GigServiceConfig{Address: "127.0.0.1:1"}, nil)
		Expect(cl).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("constructs and closes the client", func() {
		logger, err := logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())

		cl, err := NewGigClient(GigServiceConfig{Address: "127.0.0.1:1"}, logger)
		Expect(err).NotTo(HaveOccurred())
		Expect(cl).NotTo(BeNil())
		Expect(cl.Close()).To(Succeed())
		Expect((&gigClient{}).Close()).To(Succeed())
	})

	It("maps every outbound call through the fake client", func() {
		fake := &fakeGigCommandServiceClient{}
		mapr := newGigMapper(logging.Logger(nil))
		client := &gigClient{cl: fake, mapr: mapr, log: logging.Logger(nil)}

		fake.createRes = &gigv1.CreateDraftResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}
		fake.updateRes = &gigv1.UpdateBasicInfoResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}
		fake.packagesRes = &gigv1.ReplacePackagesResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}
		fake.questionsRes = &gigv1.ReplaceQuestionsResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}
		fake.mediaRes = &gigv1.ReplaceMediaResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}
		fake.draftRes = &gigv1.GetDraftResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}
		fake.slugRes = &gigv1.GetGigBySlugResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}
		fake.publishRes = &gigv1.PublishResponse{Gig: &gigv1.Gig{GigId: "gig-1"}}

		res, err := client.CreateDraft(context.Background(), gateway.CreateGigDraftRequest{FreelancerID: "freelancer-1"})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.createReq.GetFreelancerId()).To(Equal("freelancer-1"))
		Expect(res.GigID).To(Equal("gig-1"))

		res, err = client.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Title: "Title", Description: "Description", CategoryID: 1, Currency: "usd"})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.updateReq.GetGigId()).To(Equal("gig-1"))
		Expect(res.GigID).To(Equal("gig-1"))

		res, err = client.ReplacePackages(context.Background(), gateway.ReplaceGigPackagesRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Packages: []gateway.GigPackage{{Tier: gateway.TierBasic, Description: "basic", DeliveryDays: 1, PriceCents: 1}}})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.packagesReq.GetPackages()).To(HaveLen(1))
		Expect(res.GigID).To(Equal("gig-1"))

		res, err = client.ReplaceQuestions(context.Background(), gateway.ReplaceGigQuestionsRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Questions: []gateway.GigQuestion{{Content: "Question?"}}})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.questionsReq.GetQuestions()).To(HaveLen(1))
		Expect(res.GigID).To(Equal("gig-1"))

		res, err = client.ReplaceMedia(context.Background(), gateway.ReplaceGigMediaRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Files: []gateway.GigMediaUpload{{Filename: "cover.png", ContentType: "image/png", Data: []byte("x")}}})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.mediaReq.GetFiles()).To(HaveLen(1))
		Expect(res.GigID).To(Equal("gig-1"))

		res, err = client.GetDraft(context.Background(), gateway.GetGigDraftRequest{GigID: "gig-1", FreelancerID: "freelancer-1"})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.draftReq.GetGigId()).To(Equal("gig-1"))
		Expect(res.GigID).To(Equal("gig-1"))

		res, err = client.GetBySlug(context.Background(), gateway.GetGigBySlugRequest{Slug: "my-gig-gig-1"})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.slugReq.GetSlug()).To(Equal("my-gig-gig-1"))
		Expect(res.GigID).To(Equal("gig-1"))

		res, err = client.Publish(context.Background(), gateway.PublishGigRequest{GigID: "gig-1", FreelancerID: "freelancer-1"})
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.publishReq.GetGigId()).To(Equal("gig-1"))
		Expect(res.GigID).To(Equal("gig-1"))
	})

	It("maps upstream errors through the client", func() {
		fake := &fakeGigCommandServiceClient{err: status.Error(codes.NotFound, "missing")}
		client := &gigClient{cl: fake, mapr: newGigMapper(logging.Logger(nil)), log: logging.Logger(nil)}

		_, err := client.CreateDraft(context.Background(), gateway.CreateGigDraftRequest{FreelancerID: "freelancer-1"})
		Expect(err).To(MatchError(gateway.ErrGigNotFound))

		fake.err = status.Error(codes.InvalidArgument, gateway.ErrInvalidGigID.Error())
		_, err = client.UpdateBasicInfo(context.Background(), gateway.UpdateGigBasicInfoRequest{GigID: "gig-1", FreelancerID: "freelancer-1", Title: "Title", Description: "Description", CategoryID: 1, Currency: "usd"})
		Expect(err).To(MatchError(gateway.ErrInvalidGigID))
	})
})
