package http

import gateway "api-gateway/internal/domain"

var (
	_ = gateway.CreateGigDraftRequest{}
	_ = gateway.UpdateGigBasicInfoRequest{}
	_ = gateway.ReplaceGigPackagesRequest{}
	_ = gateway.ReplaceGigQuestionsRequest{}
	_ = gateway.ReplaceGigMediaRequest{}
	_ = gateway.GetGigDraftRequest{}
	_ = gateway.Gig{}
	_ = gateway.GigPreviewPage{}
	_ = gateway.GetGigBySlugRequest{}
	_ = gateway.GetMyGigsRequest{}
	_ = gateway.PublishGigRequest{}
)

// swaggerCreateGigDraftDoc documents the draft creation endpoint.
//
// @Summary Create gig draft
// @Description Creates a new freelancer gig draft owned by the authenticated user.
// @Tags gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body gateway.CreateGigDraftRequest true "Draft payload"
// @Success 202 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /gigs/drafts [post]
func swaggerCreateGigDraftDoc() {}

// swaggerUpdateGigBasicInfoDoc documents the basic-info update endpoint.
//
// @Summary Update gig basic info
// @Description Updates the public metadata of a gig draft.
// @Tags gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param gig_id path string true "Gig ID"
// @Param request body gateway.UpdateGigBasicInfoRequest true "Basic info payload"
// @Success 200 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /gigs/{gig_id}/basic-info [patch]
func swaggerUpdateGigBasicInfoDoc() {}

// swaggerReplaceGigPackagesDoc documents the package replacement endpoint.
//
// @Summary Replace gig packages
// @Description Replaces all gig pricing tiers for the authenticated draft owner.
// @Tags gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param gig_id path string true "Gig ID"
// @Param request body gateway.ReplaceGigPackagesRequest true "Packages payload"
// @Success 200 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /gigs/{gig_id}/packages [put]
func swaggerReplaceGigPackagesDoc() {}

// swaggerReplaceGigQuestionsDoc documents the requirement-question replacement endpoint.
//
// @Summary Replace gig requirements
// @Description Replaces all buyer requirement questions for a gig draft.
// @Tags gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param gig_id path string true "Gig ID"
// @Param request body gateway.ReplaceGigQuestionsRequest true "Requirements payload"
// @Success 200 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /gigs/{gig_id}/requirements [put]
func swaggerReplaceGigQuestionsDoc() {}

// swaggerReplaceGigMediaDoc documents the media replacement endpoint.
//
// @Summary Replace gig media
// @Description Replaces the ordered media batch for a gig draft.
// @Tags gigs
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param gig_id path string true "Gig ID"
// @Param files formData file true "Uploaded media files"
// @Success 200 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /gigs/{gig_id}/media [put]
func swaggerReplaceGigMediaDoc() {}

// swaggerGetGigDraftDoc documents the draft lookup endpoint.
//
// @Summary Get gig draft
// @Description Loads the current gig draft state for the authenticated owner.
// @Tags gigs
// @Produce json
// @Security BearerAuth
// @Param gig_id path string true "Gig ID"
// @Success 200 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /gigs/{gig_id}/draft [get]
func swaggerGetGigDraftDoc() {}

// swaggerGetGigBySlugDoc documents the public gig detail endpoint.
//
// @Summary Get gig by slug
// @Description Loads the public gig detail view for a seller page.
// @Tags gigs
// @Produce json
// @Param username path string true "Username"
// @Param slug path string true "Gig slug"
// @Param cursor query string false "Review cursor"
// @Success 200 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/gigs/{slug} [get]
func swaggerGetGigBySlugDoc() {}

// swaggerGetMyGigsDoc documents the authenticated owner gig list endpoint.
//
// @Summary List my gigs
// @Description Loads the authenticated freelancer's gig list and owner preview fields.
// @Tags gigs
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param status query string false "Gig status filter"
// @Param sort query string false "Gig sort key"
// @Param order query string false "Sort order"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} gateway.GigPreviewPage
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /users/{username}/gigs [get]
func swaggerGetMyGigsDoc() {}

// swaggerPublishGigDoc documents the draft publish endpoint.
//
// @Summary Publish gig
// @Description Publishes a complete gig draft.
// @Tags gigs
// @Produce json
// @Security BearerAuth
// @Param gig_id path string true "Gig ID"
// @Success 202 {object} gateway.Gig
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 412 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /gigs/{gig_id}/publish [post]
func swaggerPublishGigDoc() {}
