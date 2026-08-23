package http

import (
	"api-gateway/config"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	faultProfile500         = "status_500"
	faultProfile502         = "status_502"
	faultProfile503         = "status_503"
	faultProfileUnavailable = "unavailable"
	faultProfileTimeout     = "timeout"
	faultProfileDelay       = "delay"
)

type faultInjector struct {
	cfg         config.FaultInjectionConfig
	environment string
	startedAt   time.Time
	failures    atomic.Int64
}

func newFaultInjector(cfg config.FaultInjectionConfig, environment string, now time.Time) *faultInjector {
	return &faultInjector{cfg: cfg, environment: environment, startedAt: now}
}

func (f *faultInjector) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile, rate := f.requestConfig(c)
		if !f.shouldInject(c, profile, rate) {
			return c.Next()
		}
		c.Set("X-Fault-Injected", "true")
		if f.cfg.Delay > 0 {
			time.Sleep(f.cfg.Delay)
		}
		switch strings.ToLower(strings.TrimSpace(profile)) {
		case faultProfileDelay:
			return c.Next()
		case "mixed":
			return fiber.NewError(fiber.StatusServiceUnavailable, "fault injection")
		case faultProfile500:
			return fiber.NewError(fiber.StatusInternalServerError, "fault injection")
		case faultProfile502:
			return fiber.NewError(fiber.StatusBadGateway, "fault injection")
		case faultProfileUnavailable, faultProfile503:
			return fiber.NewError(fiber.StatusServiceUnavailable, "fault injection")
		case faultProfileTimeout:
			return fiber.NewError(fiber.StatusGatewayTimeout, "fault injection")
		default:
			return c.Next()
		}
	}
}

func (f *faultInjector) requestConfig(c *fiber.Ctx) (string, float64) {
	profile := f.cfg.Profile
	rate := f.cfg.Rate
	if f.cfg.TestToken != "" && c.Get("X-Fault-Test-Token") == f.cfg.TestToken {
		if requested := strings.TrimSpace(c.Get("X-Fault-Profile")); requested != "" {
			profile = requested
		}
		if requested := strings.TrimSpace(c.Get("X-Fault-Rate")); requested != "" {
			if parsed, err := strconv.ParseFloat(requested, 64); err == nil {
				rate = parsed
			}
		}
	}
	return profile, rate
}

func (f *faultInjector) shouldInject(c *fiber.Ctx, profile string, rate float64) bool {
	if !f.cfg.Enabled || strings.EqualFold(strings.TrimSpace(profile), "none") {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(f.environment), "production") || strings.EqualFold(strings.TrimSpace(f.environment), "prod") {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(c.Get("X-Recovery-Replay")), "true") {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(c.Get("X-Recovery-Loop-Guard")), "true") {
		return false
	}
	if f.cfg.TestToken == "" || c.Get("X-Fault-Test-Token") != f.cfg.TestToken {
		return false
	}
	requestTarget := c.Get("X-Fault-Target")
	if f.cfg.Target != "" && f.cfg.Target != "*" && requestTarget != f.cfg.Target {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(requestTarget), "gig-service") && !strings.HasPrefix(c.Path(), "/api/v2/gigs") {
		return false
	}
	if f.cfg.TTL > 0 && time.Since(f.startedAt) >= f.cfg.TTL {
		return false
	}
	if f.cfg.MaxFailures > 0 && f.failures.Load() >= f.cfg.MaxFailures {
		return false
	}
	if rate > 1 {
		rate /= 100
	}
	if rate > 0 && rate < 1 && rand.Float64() >= rate {
		return false
	}
	f.failures.Add(1)
	return true
}
