package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"github.com/google/uuid"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type gigHandler struct {
	service GigService
	auth    GigPrincipalResolver
	log     logging.Logger
}

// NewGigHandler constructs the gig HTTP handler group.
func NewGigHandler(service GigService, jwtSecret string, log logging.Logger) (GigHandler, error) {
	if service == nil {
		return nil, ErrNilGigService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	auth, err := newJWTPrincipalResolver(jwtSecret, log)
	if err != nil {
		return nil, err
	}

	return &gigHandler{
		service: service,
		auth:    auth,
		log:     log.With(logging.String("module", "http-gig-handler")),
	}, nil
}

// RegisterRoutes mounts gig routes under the router it receives.
func (h *gigHandler) RegisterRoutes(router fiber.Router) {
	ownerGigs := router.Group("/users/:username/gigs")
	ownerGigs.Get("", h.auth.Middleware(), h.HandleGetMyGigs)

	publicGigs := router.Group("/users/:username/gigs")
	publicGigs.Get("/*", h.HandleGetBySlug)

	authGigs := router.Group("/gigs")
	authGigs.Use(h.auth.Middleware())
	authGigs.Post("/drafts", h.HandleCreateDraft)
	authGigs.Patch("/:gig_id/basic-info", h.HandleUpdateBasicInfo)
	authGigs.Put("/:gig_id/packages", h.HandleReplacePackages)
	authGigs.Put("/:gig_id/requirements", h.HandleReplaceQuestions)
	authGigs.Put("/:gig_id/media", h.HandleReplaceMedia)
	authGigs.Get("/:gig_id/draft", h.HandleGetDraft)
	authGigs.Post("/:gig_id/publish", h.HandlePublish)
}

// HandleCreateDraft starts a new gig draft.
func (h *gigHandler) HandleCreateDraft(c *fiber.Ctx) error {
	var req gateway.CreateGigDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.FreelancerID = freelancerID

	result, err := h.service.CreateDraft(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	h.log.Info("gig draft created")
	return c.Status(fiber.StatusAccepted).JSON(result)
}

// HandleUpdateBasicInfo updates the draft's public metadata.
func (h *gigHandler) HandleUpdateBasicInfo(c *fiber.Ctx) error {
	var req gateway.UpdateGigBasicInfoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.GigID = c.Params("gig_id")
	req.FreelancerID = freelancerID

	result, err := h.service.UpdateBasicInfo(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleReplacePackages replaces the gig's pricing tiers.
func (h *gigHandler) HandleReplacePackages(c *fiber.Ctx) error {
	var req gateway.ReplaceGigPackagesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.GigID = c.Params("gig_id")
	req.FreelancerID = freelancerID

	result, err := h.service.ReplacePackages(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleReplaceQuestions replaces the gig requirements questions.
func (h *gigHandler) HandleReplaceQuestions(c *fiber.Ctx) error {
	var req gateway.ReplaceGigQuestionsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": ErrInvalidRequestBody.Error()})
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.GigID = c.Params("gig_id")
	req.FreelancerID = freelancerID

	result, err := h.service.ReplaceQuestions(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleReplaceMedia replaces the gig media references.
func (h *gigHandler) HandleReplaceMedia(c *fiber.Ctx) error {
	files, err := h.readMediaUploads(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req := gateway.ReplaceGigMediaRequest{
		GigID:        c.Params("gig_id"),
		FreelancerID: freelancerID,
		Files:        files,
	}

	result, err := h.service.ReplaceMedia(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleGetDraft loads the current draft state.
func (h *gigHandler) HandleGetDraft(c *fiber.Ctx) error {
	req := gateway.GetGigDraftRequest{
		GigID: c.Params("gig_id"),
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.FreelancerID = freelancerID

	result, err := h.service.GetDraft(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleGetBySlug loads the public gig detail view by slug.
func (h *gigHandler) HandleGetBySlug(c *fiber.Ctx) error {
	path := strings.TrimSpace(c.Params("*"))
	if path == "" {
		return h.mapGigError(c, gateway.ErrGigNotFound)
	}
	if !strings.Contains(path, "-") || !hasValidGigIDSuffix(path) {
		return h.mapGigError(c, gateway.ErrGigNotFound)
	}
	result, err := h.service.GetBySlug(c.UserContext(), gateway.GetGigBySlugRequest{
		Username: c.Params("username"),
		Slug:     path,
		Cursor:   c.Query("cursor"),
	})
	if err != nil {
		return h.mapGigError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// HandleGetMyGigs loads the authenticated owner's gig list.
func (h *gigHandler) HandleGetMyGigs(c *fiber.Ctx) error {
	username, err := url.PathUnescape(c.Params("username"))
	if err != nil {
		return h.mapGigError(c, gateway.ErrInvalidUsername)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return h.mapGigError(c, gateway.ErrInvalidUsername)
	}
	actorUsername, err := h.auth.Username(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	if strings.TrimSpace(actorUsername) != username {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	page, err := parsePositiveIntQuery(c.Query("page"), 1, gateway.ErrInvalidGigListPage)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	limit, err := parsePositiveIntQuery(c.Query("limit"), 10, gateway.ErrInvalidGigListLimit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req := gateway.GetMyGigsRequest{
		UserID: freelancerID,
		Status: strings.TrimSpace(c.Query("status")),
		Sort:   strings.TrimSpace(c.Query("sort")),
		Order:  strings.TrimSpace(c.Query("order")),
		Page:   int32(page),
		Limit:  int32(limit),
	}

	result, err := h.service.GetMyGigs(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	h.log.Info("gig list request accepted",
		logging.Operation("http.gig.list"),
		logging.String("username", username),
		logging.String("user_id", freelancerID),
		logging.String("status", req.Status),
		logging.String("sort", req.Sort),
		logging.String("order", req.Order),
	)
	return c.Status(fiber.StatusOK).JSON(result)
}

// HandlePublish publishes a complete gig draft.
func (h *gigHandler) HandlePublish(c *fiber.Ctx) error {
	req := gateway.PublishGigRequest{
		GigID: c.Params("gig_id"),
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	username, err := h.auth.Username(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.FreelancerID = freelancerID
	req.Username = username

	result, err := h.service.Publish(c.UserContext(), req)
	if err != nil {
		return h.mapGigError(c, err)
	}

	return c.Status(fiber.StatusAccepted).JSON(result)
}

func (h *gigHandler) mapGigError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, gateway.ErrInvalidGigID),
		errors.Is(err, gateway.ErrInvalidFreelancerID),
		errors.Is(err, gateway.ErrInvalidTitle),
		errors.Is(err, gateway.ErrInvalidDescription),
		errors.Is(err, gateway.ErrInvalidCategoryID),
		errors.Is(err, gateway.ErrInvalidCurrency),
		errors.Is(err, gateway.ErrInvalidPackageTier),
		errors.Is(err, gateway.ErrInvalidPackageDescription),
		errors.Is(err, gateway.ErrInvalidPackageDeliveryDays),
		errors.Is(err, gateway.ErrInvalidPackagePriceCents),
		errors.Is(err, gateway.ErrInvalidQuestionContent),
		errors.Is(err, gateway.ErrInvalidMediaUpload),
		errors.Is(err, gateway.ErrInvalidPackageCount):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrGigNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrInvalidGigListStatus),
		errors.Is(err, gateway.ErrInvalidGigListSort),
		errors.Is(err, gateway.ErrInvalidGigListOrder),
		errors.Is(err, gateway.ErrInvalidGigListPage),
		errors.Is(err, gateway.ErrInvalidGigListLimit),
		errors.Is(err, gateway.ErrInvalidUserID):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrGigDraftIncomplete),
		errors.Is(err, gateway.ErrGigAlreadyPublished),
		errors.Is(err, gateway.ErrConnectOnboardingIncomplete),
		errors.Is(err, gateway.ErrInvalidGigState):
		return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrInvalidGigSlug):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrInvalidUsername):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, gateway.ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	default:
		h.log.Error("request failed", logging.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

func parsePositiveIntQuery(raw string, def int, invalidErr error) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return def, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, invalidErr
	}
	return value, nil
}

func hasValidGigIDSuffix(path string) bool {
	if len(path) < 36 {
		return false
	}
	_, err := uuid.Parse(path[len(path)-36:])
	return err == nil
}

func (h *gigHandler) readMediaUploads(c *fiber.Ctx) ([]gateway.GigMediaUpload, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, ErrInvalidRequestBody
	}

	headers := form.File["files"]
	if len(headers) == 0 {
		return nil, gateway.ErrInvalidMediaUpload
	}

	uploads := make([]gateway.GigMediaUpload, 0, len(headers))
	for _, header := range headers {
		upload, err := readMultipartFile(header)
		if err != nil {
			return nil, ErrInvalidRequestBody
		}
		uploads = append(uploads, upload)
	}

	return uploads, nil
}

func readMultipartFile(header *multipart.FileHeader) (gateway.GigMediaUpload, error) {
	file, err := header.Open()
	if err != nil {
		return gateway.GigMediaUpload{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return gateway.GigMediaUpload{}, err
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return gateway.GigMediaUpload{
		Filename:    header.Filename,
		ContentType: contentType,
		Data:        data,
	}, nil
}
