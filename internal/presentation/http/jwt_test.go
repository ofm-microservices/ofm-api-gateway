package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ofm-microservices/ofm-common/pkg/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/valyala/fasthttp"
)

var _ = Describe("jwt principal resolver", func() {
	var logger logging.Logger

	BeforeEach(func() {
		var err error
		logger, err = logging.New("api-gateway", "test", "debug")
		Expect(err).NotTo(HaveOccurred())
	})

	It("validates constructor inputs", func() {
		resolver, err := newJWTPrincipalResolver("", logger)
		Expect(resolver).To(BeNil())
		Expect(err).To(MatchError(errEmptyJWTSecret))

		resolver, err = newJWTPrincipalResolver("secret", nil)
		Expect(resolver).To(BeNil())
		Expect(err).To(MatchError(ErrNilLogger))
	})

	It("parses bearer headers and tokens", func() {
		_, err := bearerToken("")
		Expect(err).To(MatchError(errInvalidAuthorizationHeader))
		_, err = bearerToken("Token abc")
		Expect(err).To(MatchError(errInvalidAuthorizationHeader))
		token, err := bearerToken("Bearer abc")
		Expect(err).NotTo(HaveOccurred())
		Expect(token).To(Equal("abc"))

		_, _, _, err = splitJWT("a.b")
		Expect(err).To(MatchError(errInvalidJWTToken))
		h, p, s, err := splitJWT("a.b.c")
		Expect(err).NotTo(HaveOccurred())
		Expect([]string{h, p, s}).To(Equal([]string{"a", "b", "c"}))
	})

	It("extracts freelancer ids from signed jwt tokens", func() {
		resolver, err := newJWTPrincipalResolver("secret", logger)
		Expect(err).NotTo(HaveOccurred())

		token := signedJWT("freelancer-1", "secret")
		id, err := resolver.extractFreelancerID("Bearer " + token)
		Expect(err).NotTo(HaveOccurred())
		Expect(id).To(Equal("freelancer-1"))

		_, err = resolver.extractFreelancerID("Bearer " + signedJWT("freelancer-1", "different"))
		Expect(err).To(MatchError(errInvalidJWTToken))

		expired := signedJWTWithExp("freelancer-1", "secret", time.Now().Add(-time.Hour))
		_, err = resolver.extractFreelancerID("Bearer " + expired)
		Expect(err).To(MatchError(errExpiredJWTToken))

		_, err = resolver.extractFreelancerID("Bearer " + signedJWTWithoutSubject("secret"))
		Expect(err).To(MatchError(errInvalidJWTToken))
	})

	It("populates middleware locals and reads freelancer ids", func() {
		resolver, err := newJWTPrincipalResolver("secret", logger)
		Expect(err).NotTo(HaveOccurred())

		app := fiber.New()
		app.Use(resolver.Middleware())
		app.Get("/who", func(c *fiber.Ctx) error {
			id, err := resolver.FreelancerID(c)
			Expect(err).NotTo(HaveOccurred())
			return c.SendString(id)
		})

		req, _ := http.NewRequest(http.MethodGet, "/who", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+signedJWT("freelancer-1", "secret"))
		resp, err := app.Test(req, -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(fiber.StatusOK))

		c := app.AcquireCtx(&fasthttp.RequestCtx{})
		_, err = resolver.FreelancerID(c)
		Expect(err).To(MatchError(errMissingJWTPrincipal))
	})
})

func signedJWTWithoutSubject(secret string) string {
	return signedJWTWithExp("", secret, time.Now().Add(time.Hour))
}

func signedJWTWithExp(subject, secret string, exp time.Time) string {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := map[string]any{
		"iat": time.Now().UTC().Unix(),
		"exp": exp.UTC().Unix(),
	}
	if subject != "" {
		claims["sub"] = subject
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)
	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsPart := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := headerPart + "." + claimsPart

	mac := hmacForTest(secret, signingInput)
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(mac)
}

func hmacForTest(secret, signingInput string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signingInput))
	return mac.Sum(nil)
}
