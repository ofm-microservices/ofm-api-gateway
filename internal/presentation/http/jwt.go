package http

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	commonjwt "github.com/ofm-microservices/ofm-common/pkg/jwt"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
)

var (
	errEmptyJWTSecret             = errors.New("jwt secret is empty")
	errInvalidAuthorizationHeader = errors.New("invalid authorization header")
	errInvalidJWTToken            = errors.New("invalid jwt token")
	errExpiredJWTToken            = errors.New("expired jwt token")
	errMissingJWTPrincipal        = errors.New("missing jwt principal")
)

const jwtPrincipalLocalKey = "jwt.freelancer_id"
const jwtClaimsLocalKey = "jwt.claims"

type jwtPrincipalResolver struct {
	verifier commonjwt.Verifier
	log      logging.Logger
}

func newJWTPrincipalResolver(secret string, log logging.Logger) (*jwtPrincipalResolver, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, errEmptyJWTSecret
	}
	if log == nil {
		return nil, ErrNilLogger
	}
	verifier, err := commonjwt.NewVerifier(commonjwt.Config{Secret: secret})
	if err != nil {
		return nil, err
	}

	return &jwtPrincipalResolver{
		verifier: verifier,
		log:      log.With(logging.String("module", "jwt-principal-resolver")),
	}, nil
}

func (r *jwtPrincipalResolver) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("X-Recovery-Replay") == "true" && strings.TrimSpace(c.Get("Authorization")) == "" {
			if principal := strings.TrimSpace(c.Get("X-Recovery-Principal-ID")); principal != "" {
				c.Locals(jwtPrincipalLocalKey, principal)
				// Recovery commands deliberately do not carry a user's bearer
				// token through Kafka. Populate the same claims boundary used by
				// handlers that only need actor identity; domain authorization
				// still receives the explicit recovery principal.
				c.Locals(jwtClaimsLocalKey, &commonjwt.Claims{
					Subject:  principal,
					Username: strings.TrimSpace(c.Get("X-Recovery-Username")),
					Email:    strings.TrimSpace(c.Get("X-Recovery-Email")),
				})
				return c.Next()
			}
		}
		claims, err := r.extractClaims(c.Get("Authorization"))
		if err != nil {
			r.log.Error("jwt authorization failed", logging.Err(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		c.Locals(jwtPrincipalLocalKey, strings.TrimSpace(claims.Subject))
		c.Locals(jwtClaimsLocalKey, claims)
		return c.Next()
	}
}

func (r *jwtPrincipalResolver) FreelancerID(c *fiber.Ctx) (string, error) {
	value, ok := c.Locals(jwtPrincipalLocalKey).(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", errMissingJWTPrincipal
	}

	return strings.TrimSpace(value), nil
}

func (r *jwtPrincipalResolver) extractClaims(header string) (*commonjwt.Claims, error) {
	token, err := commonjwt.ParseBearer(header)
	if err != nil {
		switch {
		case errors.Is(err, commonjwt.ErrInvalidAuthorizationHeader):
			return nil, errInvalidAuthorizationHeader
		default:
			return nil, errInvalidJWTToken
		}
	}
	claims, err := r.verifier.Validate(token)
	if err != nil {
		switch {
		case errors.Is(err, commonjwt.ErrExpiredToken):
			return nil, fmt.Errorf("%w: %v", errExpiredJWTToken, err)
		case errors.Is(err, commonjwt.ErrInvalidAuthorizationHeader):
			return nil, fmt.Errorf("%w: %v", errInvalidAuthorizationHeader, err)
		default:
			return nil, fmt.Errorf("%w: %v", errInvalidJWTToken, err)
		}
	}
	if strings.TrimSpace(claims.Subject) == "" {
		return nil, errInvalidJWTToken
	}
	return claims, nil
}

func (r *jwtPrincipalResolver) Claims(c *fiber.Ctx) (*commonjwt.Claims, error) {
	value, ok := c.Locals(jwtClaimsLocalKey).(*commonjwt.Claims)
	if !ok || value == nil {
		return nil, errMissingJWTPrincipal
	}
	return value, nil
}

func (r *jwtPrincipalResolver) Username(c *fiber.Ctx) (string, error) {
	claims, err := r.Claims(c)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(claims.Username) == "" {
		return "", errMissingJWTPrincipal
	}
	return strings.TrimSpace(claims.Username), nil
}

func (r *jwtPrincipalResolver) Email(c *fiber.Ctx) (string, error) {
	claims, err := r.Claims(c)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(claims.Email) == "" {
		return "", errMissingJWTPrincipal
	}
	return strings.TrimSpace(claims.Email), nil
}

func (r *jwtPrincipalResolver) Roles(c *fiber.Ctx) ([]string, error) {
	claims, err := r.Claims(c)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0, len(claims.Roles))
	for _, role := range claims.Roles {
		if strings.TrimSpace(role) == "" {
			continue
		}
		roles = append(roles, strings.TrimSpace(role))
	}
	return roles, nil
}

func (r *jwtPrincipalResolver) HasRole(c *fiber.Ctx, role string) (bool, error) {
	roles, err := r.Roles(c)
	if err != nil {
		return false, err
	}
	role = strings.TrimSpace(role)
	if role == "" {
		return false, nil
	}
	for _, candidate := range roles {
		if strings.EqualFold(candidate, role) {
			return true, nil
		}
	}
	return false, nil
}

func bearerToken(header string) (string, error) {
	return commonjwt.ParseBearer(header)
}

func splitJWT(token string) (string, string, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", "", "", errInvalidJWTToken
	}
	return parts[0], parts[1], parts[2], nil
}
