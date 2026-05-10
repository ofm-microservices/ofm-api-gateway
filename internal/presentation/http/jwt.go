package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microseervices/ofm-common/pkg/logging"
)

var (
	errEmptyJWTSecret             = errors.New("jwt secret is empty")
	errInvalidAuthorizationHeader = errors.New("invalid authorization header")
	errInvalidJWTToken            = errors.New("invalid jwt token")
	errExpiredJWTToken            = errors.New("expired jwt token")
	errMissingJWTPrincipal        = errors.New("missing jwt principal")
)

const jwtPrincipalLocalKey = "jwt.freelancer_id"

type jwtPrincipalResolver struct {
	secret []byte
	log    logging.Logger
}

func newJWTPrincipalResolver(secret string, log logging.Logger) (*jwtPrincipalResolver, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, errEmptyJWTSecret
	}
	if log == nil {
		return nil, ErrNilLogger
	}

	return &jwtPrincipalResolver{
		secret: []byte(secret),
		log:    log.With(logging.String("module", "jwt-principal-resolver")),
	}, nil
}

func (r *jwtPrincipalResolver) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		freelancerID, err := r.extractFreelancerID(c.Get("Authorization"))
		if err != nil {
			r.log.Error("jwt authorization failed", logging.Err(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		c.Locals(jwtPrincipalLocalKey, freelancerID)
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

func (r *jwtPrincipalResolver) extractFreelancerID(header string) (string, error) {
	token, err := bearerToken(header)
	if err != nil {
		return "", err
	}

	headerPart, payloadPart, signaturePart, err := splitJWT(token)
	if err != nil {
		return "", err
	}

	signed := headerPart + "." + payloadPart
	expected := hmac.New(sha256.New, r.secret)
	if _, err := expected.Write([]byte(signed)); err != nil {
		return "", errInvalidJWTToken
	}

	signature, err := base64.RawURLEncoding.DecodeString(signaturePart)
	if err != nil {
		return "", errInvalidJWTToken
	}
	if !hmac.Equal(signature, expected.Sum(nil)) {
		return "", errInvalidJWTToken
	}

	var claims map[string]any
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return "", errInvalidJWTToken
	}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", errInvalidJWTToken
	}

	if exp, ok := claims["exp"].(float64); ok && time.Now().UTC().After(time.Unix(int64(exp), 0)) {
		return "", errExpiredJWTToken
	}

	sub, _ := claims["sub"].(string)
	if strings.TrimSpace(sub) == "" {
		return "", errInvalidJWTToken
	}

	return strings.TrimSpace(sub), nil
}

func bearerToken(header string) (string, error) {
	fields := strings.Fields(header)
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return "", errInvalidAuthorizationHeader
	}

	return fields[1], nil
}

func splitJWT(token string) (string, string, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", "", "", errInvalidJWTToken
	}

	return parts[0], parts[1], parts[2], nil
}
