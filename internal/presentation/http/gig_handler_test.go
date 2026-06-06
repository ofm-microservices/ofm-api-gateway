package http

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	gateway "api-gateway/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type gigServiceStub struct {
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

	err error
}

func (s *gigServiceStub) CreateDraft(_ context.Context, req gateway.CreateGigDraftRequest) (*gateway.Gig, error) {
	s.createReq = req
	return s.createRes, s.err
}

func (s *gigServiceStub) UpdateBasicInfo(_ context.Context, req gateway.UpdateGigBasicInfoRequest) (*gateway.Gig, error) {
	s.updateReq = req
	return s.updateRes, s.err
}

func (s *gigServiceStub) ReplacePackages(_ context.Context, req gateway.ReplaceGigPackagesRequest) (*gateway.Gig, error) {
	s.packagesReq = req
	return s.packagesRes, s.err
}

func (s *gigServiceStub) ReplaceQuestions(_ context.Context, req gateway.ReplaceGigQuestionsRequest) (*gateway.Gig, error) {
	s.questionsReq = req
	return s.questionsRes, s.err
}

func (s *gigServiceStub) ReplaceMedia(_ context.Context, req gateway.ReplaceGigMediaRequest) (*gateway.Gig, error) {
	s.mediaReq = req
	return s.mediaRes, s.err
}

func (s *gigServiceStub) GetDraft(_ context.Context, req gateway.GetGigDraftRequest) (*gateway.Gig, error) {
	s.draftReq = req
	return s.draftRes, s.err
}

func (s *gigServiceStub) GetBySlug(_ context.Context, req gateway.GetGigBySlugRequest) (*gateway.Gig, error) {
	s.slugReq = req
	return s.slugRes, s.err
}

func (s *gigServiceStub) GetPreviewGigsByFreelancerUsername(_ context.Context, req gateway.GetPreviewGigsByFreelancerUsernameRequest) (*gateway.GigPreviewList, error) {
	s.previewReq = req
	return s.previewRes, s.err
}

func (s *gigServiceStub) GetMyGigs(_ context.Context, req gateway.GetMyGigsRequest) (*gateway.GigPreviewPage, error) {
	s.myGigsReq = req
	return s.myGigsRes, s.err
}

func (s *gigServiceStub) Publish(_ context.Context, req gateway.PublishGigRequest) (*gateway.Gig, error) {
	s.publishReq = req
	return s.publishRes, s.err
}

var _ = Describe("GigHandler", func() {
	var (
		service *gigServiceStub
		logger  logging.Logger
		app     *fiber.App
	)

	const testJWTSecret = "local-dev-secret-change-me"

	BeforeEach(func() {
		service = &gigServiceStub{}
		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())

		handler, err := NewGigHandler(service, testJWTSecret, logger)
		Expect(err).NotTo(HaveOccurred())
		app = fiber.New()
		v1 := app.Group("/v1")
		handler.RegisterRoutes(v1)
	})

	It("validates constructor input", func() {
		handler, err := NewGigHandler(nil, testJWTSecret, logger)
		Expect(handler).To(BeNil())
		Expect(err).To(MatchError(ErrNilGigService))

		handler, err = NewGigHandler(service, testJWTSecret, nil)
		Expect(handler).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))

		handler, err = NewGigHandler(service, " ", logger)
		Expect(handler).To(BeNil())
		Expect(err).To(HaveOccurred())
	})

	It("creates drafts", func() {
		service.createRes = &gateway.Gig{GigID: "gig-1"}
		resp, err := app.Test(gigJSONRequest("POST", "/v1/gigs/drafts", `{}`, signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusAccepted))
		Expect(service.createReq.FreelancerID).To(Equal("freelancer-1"))
	})

	It("rejects invalid create payloads", func() {
		resp, err := app.Test(gigJSONRequest("POST", "/v1/gigs/drafts", "{", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
	})

	It("updates basic info", func() {
		service.updateRes = &gateway.Gig{GigID: "gig-1"}
		resp, err := app.Test(gigJSONRequest("PATCH", "/v1/gigs/gig-1/basic-info", `{"title":"Title","description":"Description","category_id":1,"currency":"usd"}`, signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(service.updateReq.GigID).To(Equal("gig-1"))
		Expect(service.updateReq.FreelancerID).To(Equal("freelancer-1"))
	})

	It("replaces packages, questions, media, loads drafts, and publishes", func() {
		service.packagesRes = &gateway.Gig{GigID: "gig-1"}
		service.questionsRes = &gateway.Gig{GigID: "gig-1"}
		service.mediaRes = &gateway.Gig{GigID: "gig-1"}
		service.draftRes = &gateway.Gig{GigID: "gig-1"}
		service.publishRes = &gateway.Gig{GigID: "gig-1"}

		resp, err := app.Test(gigJSONRequest("PUT", "/v1/gigs/gig-1/packages", `{"packages":[{"tier":"basic","description":"Basic","delivery_days":1,"price_cents":1000},{"tier":"standard","description":"Standard","delivery_days":2,"price_cents":2000}]}`, signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(service.packagesReq.Packages).To(HaveLen(2))

		resp, err = app.Test(gigJSONRequest("PUT", "/v1/gigs/gig-1/requirements", `{"questions":[{"content":"Question?"}]}`, signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(service.questionsReq.Questions).To(HaveLen(1))

		resp, err = app.Test(gigMultipartRequest("/v1/gigs/gig-1/media", signedJWT("freelancer-1", testJWTSecret), map[string]string{"files": "cover.png"}), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(service.mediaReq.Files).To(HaveLen(1))
		Expect(service.mediaReq.Files[0].ContentType).To(Equal("application/octet-stream"))

		resp, err = app.Test(gigJSONRequest("GET", "/v1/gigs/gig-1/draft", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))

		resp, err = app.Test(gigJSONRequest("POST", "/v1/gigs/gig-1/publish", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusAccepted))
	})

	It("rejects missing media uploads", func() {
		resp, err := app.Test(gigMultipartRequest("/v1/gigs/gig-1/media", signedJWT("freelancer-1", testJWTSecret), nil), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
	})

	It("loads drafts and publishes gigs", func() {
		service.draftRes = &gateway.Gig{GigID: "gig-1"}
		service.publishRes = &gateway.Gig{GigID: "gig-1"}

		resp, err := app.Test(gigJSONRequest("GET", "/v1/gigs/gig-1/draft", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))

		resp, err = app.Test(gigJSONRequest("POST", "/v1/gigs/gig-1/publish", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusAccepted))
		Expect(service.publishReq.FreelancerID).To(Equal("freelancer-1"))
		Expect(service.publishReq.Username).To(Equal("freelancer-1"))
	})

	It("loads a public gig by slug without auth", func() {
		service.slugRes = &gateway.Gig{GigID: "019e706c-616e-7473-9c1a-838c33b75013", Slug: "my-gig-019e706c-616e-7473-9c1a-838c33b75013"}
		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/alex/gigs/my-gig-019e706c-616e-7473-9c1a-838c33b75013?cursor=abc", "", ""), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(service.slugReq.Username).To(Equal("alex"))
		Expect(service.slugReq.Slug).To(Equal("my-gig-019e706c-616e-7473-9c1a-838c33b75013"))
		Expect(service.slugReq.Cursor).To(Equal("abc"))
	})

	It("loads owner gigs by username with auth", func() {
		service.myGigsRes = &gateway.GigPreviewPage{Items: []gateway.GigPreview{{GigID: "gig-1"}}, Page: 1, Limit: 10, TotalPages: 1}
		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/freelancer-1/gigs?status=published&sort=updated_at&order=desc&page=2&limit=10", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(service.myGigsReq.UserID).To(Equal("freelancer-1"))
		Expect(service.myGigsReq.Status).To(Equal("published"))
		Expect(service.myGigsReq.Sort).To(Equal("updated_at"))
		Expect(service.myGigsReq.Order).To(Equal("desc"))
		Expect(service.myGigsReq.Page).To(Equal(int32(2)))
		Expect(service.myGigsReq.Limit).To(Equal(int32(10)))
	})

	It("rejects owner gigs when username does not match token", func() {
		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/alex/gigs?page=1&limit=10", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusForbidden))
	})

	It("returns not found for malformed public gig slugs", func() {
		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/alex/gigs/ss", "", ""), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusNotFound))
	})

	It("rejects requests without a bearer token", func() {
		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/alex/gigs?page=1&limit=10", "", ""), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
	})

	It("maps validation errors to a bad request", func() {
		service.err = gateway.ErrInvalidGigID
		resp, err := app.Test(gigJSONRequest("GET", "/v1/gigs/gig-1/draft", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
		body, _ := io.ReadAll(resp.Body)
		Expect(string(body)).To(ContainSubstring("invalid gig id"))
	})

	It("maps not found and precondition errors", func() {
		service.err = gateway.ErrGigNotFound
		resp, err := app.Test(gigJSONRequest("GET", "/v1/gigs/gig-1/draft", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusNotFound))

		service.err = gateway.ErrGigDraftIncomplete
		resp, err = app.Test(gigJSONRequest("POST", "/v1/gigs/gig-1/publish", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusPreconditionFailed))

		service.err = gateway.ErrGigAlreadyPublished
		resp, err = app.Test(gigJSONRequest("POST", "/v1/gigs/gig-1/publish", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusPreconditionFailed))

		service.err = gateway.ErrConnectOnboardingIncomplete
		resp, err = app.Test(gigJSONRequest("POST", "/v1/gigs/gig-1/publish", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusPreconditionFailed))

		service.err = gateway.ErrInvalidGigState
		resp, err = app.Test(gigJSONRequest("GET", "/v1/gigs/gig-1/draft", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusPreconditionFailed))
	})

	It("maps unexpected errors to 500", func() {
		service.err = errors.New("boom")
		resp, err := app.Test(gigJSONRequest("GET", "/v1/gigs/gig-1/draft", "", signedJWT("freelancer-1", testJWTSecret)), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusInternalServerError))
	})
})

func gigMultipartRequest(target, token string, files map[string]string) *http.Request {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for field, name := range files {
		part, err := writer.CreateFormFile(field, name)
		Expect(err).NotTo(HaveOccurred())
		_, err = part.Write([]byte("file-data"))
		Expect(err).NotTo(HaveOccurred())
	}
	Expect(writer.Close()).To(Succeed())

	req, _ := http.NewRequest("PUT", target, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func gigJSONRequest(method, target, body, token string) *http.Request {
	req, _ := http.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func signedJWT(subject, secret string) string {
	now := time.Now().UTC()
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := map[string]any{
		"sub":      subject,
		"username": subject,
		"iat":      now.Unix(),
		"exp":      now.Add(time.Hour).Unix(),
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsPart := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := headerPart + "." + claimsPart

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	signaturePart := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signaturePart
}
