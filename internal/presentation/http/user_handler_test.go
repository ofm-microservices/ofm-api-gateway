package http

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	gateway "api-gateway/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type userProfileServiceStub struct {
	req gateway.GetUserProfileRequest
	res *gateway.UserProfile
	err error
}

func (s *userProfileServiceStub) GetUserProfile(_ context.Context, req gateway.GetUserProfileRequest) (*gateway.UserProfile, error) {
	s.req = req
	return s.res, s.err
}

var _ = Describe("UserHandler", func() {
	var (
		service *userProfileServiceStub
		logger  logging.Logger
		app     *fiber.App
	)

	BeforeEach(func() {
		service = &userProfileServiceStub{}
		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())

		handler, err := NewUserHandler(service, logger)
		Expect(err).NotTo(HaveOccurred())
		app = fiber.New()
		v1 := app.Group("/v1")
		handler.RegisterRoutes(v1)
	})

	It("validates constructor input", func() {
		handler, err := NewUserHandler(nil, logger)
		Expect(handler).To(BeNil())
		Expect(err).To(MatchError(ErrNilUserProfileService))

		handler, err = NewUserHandler(service, nil)
		Expect(handler).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("loads the composite profile page with independent cursors", func() {
		service.res = &gateway.UserProfile{
			User: &gateway.User{UserID: "user-1", Username: "alex", DisplayName: "Alex Tester"},
			Gigs: &gateway.GigPreviewList{
				Items:   []gateway.GigPreview{{GigID: "gig-1"}},
				Cursor:  "gigs-cursor",
				HasMore: true,
			},
			Reviews: &gateway.ReviewList{
				Items:   []gateway.Review{{ReviewID: "review-1"}},
				Cursor:  "reviews-cursor",
				HasMore: false,
			},
		}

		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/alex?gigs_cursor=abc&reviews_cursor=def", "", ""), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		Expect(service.req).To(Equal(gateway.GetUserProfileRequest{
			Username:      "alex",
			GigsCursor:    "abc",
			ReviewsCursor: "def",
		}))

		body, _ := io.ReadAll(resp.Body)
		var decoded map[string]any
		Expect(json.Unmarshal(body, &decoded)).To(Succeed())
		Expect(decoded["user"]).NotTo(BeNil())
		Expect(decoded["gigs"]).NotTo(BeNil())
		Expect(decoded["reviews"]).NotTo(BeNil())
		Expect(strings.Contains(string(body), "\"gigs_cursor\"")).To(BeFalse())
	})

	It("maps invalid usernames to 400", func() {
		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/ ", "", ""), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusBadRequest))
	})

	It("maps missing users to 404", func() {
		service.err = gateway.ErrUserNotFound
		resp, err := app.Test(gigJSONRequest("GET", "/v1/users/alex", "", ""), -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusNotFound))
	})
})
