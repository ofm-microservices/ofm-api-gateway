package http

import (
	gateway "api-gateway/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	"time"
)

type searchHandler struct {
	service SearchService
	log     logging.Logger
}

// NewSearchHandler constructs the public search HTTP handler group.
func NewSearchHandler(service SearchService, log logging.Logger) (SearchHandler, error) {
	if service == nil {
		return nil, ErrNilSearchService
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &searchHandler{
		service: service,
		log:     log.With(logging.String("module", "http-search-handler")),
	}, nil
}

func (h *searchHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/search", h.HandleSearch)
}

func (h *searchHandler) HandleSearch(c *fiber.Ctx) error {
	started := time.Now()
	req := gateway.SearchRequest{
		Query:  c.Query("query"),
		Cursor: c.Query("cursor"),
	}
	if v := c.Query("sort"); v != "" {
		req.Sort = parseInt32(v)
	}
	if v := c.Query("order"); v != "" {
		req.Order = parseInt32(v)
	}

	res, err := h.service.Search(c.UserContext(), req)
	if err != nil {
		h.log.Error("search request failed", logging.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	h.log.Info("search request accepted",
		logging.Operation("http.search"),
		logging.DurationMS(time.Since(started)),
	)
	return c.Status(fiber.StatusOK).JSON(res)
}

func parseInt32(value string) int32 {
	var n int32
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + int32(ch-'0')
	}
	return n
}
