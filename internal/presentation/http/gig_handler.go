package http

import (
	gateway "api-gateway/internal/domain"
	"errors"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
	"io"
	"mime/multipart"

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
	gigs := router.Group("/gigs")
	gigs.Use(h.auth.Middleware())

	gigs.Post("/drafts", h.HandleCreateDraft)
	gigs.Patch("/:gig_id/basic-info", h.HandleUpdateBasicInfo)
	gigs.Put("/:gig_id/packages", h.HandleReplacePackages)
	gigs.Put("/:gig_id/requirements", h.HandleReplaceQuestions)
	gigs.Put("/:gig_id/media", h.HandleReplaceMedia)
	gigs.Get("/:gig_id/draft", h.HandleGetDraft)
	gigs.Post("/:gig_id/publish", h.HandlePublish)
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

// HandlePublish publishes a complete gig draft.
func (h *gigHandler) HandlePublish(c *fiber.Ctx) error {
	req := gateway.PublishGigRequest{
		GigID: c.Params("gig_id"),
	}
	freelancerID, err := h.auth.FreelancerID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	req.FreelancerID = freelancerID

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
	case errors.Is(err, gateway.ErrGigDraftIncomplete),
		errors.Is(err, gateway.ErrGigAlreadyPublished),
		errors.Is(err, gateway.ErrInvalidGigState):
		return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{"error": err.Error()})
	default:
		h.log.Error("request failed", logging.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
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
