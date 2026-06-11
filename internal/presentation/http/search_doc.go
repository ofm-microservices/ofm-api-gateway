package http

import gateway "api-gateway/internal/domain"

var _ = gateway.SearchResponse{}

// swaggerSearchDoc documents the public search endpoint.
//
// @Summary Search
// @Description Searches the public gig and freelancer catalog.
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Param cursor query string false "Pagination cursor"
// @Param sort query int false "Sort key"
// @Param order query int false "Sort order"
// @Success 200 {object} gateway.SearchResult
// @Failure 400 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /search [get]
func swaggerSearchDoc() {}
