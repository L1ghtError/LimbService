package middleware

import (
	"light-backend/service"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

const HeaderTokenLookup string = "header:" + fiber.HeaderAuthorization
const HeaderAuthScheme = "Bearer"

const CookieJWT string = "jwt"
const CookieTokenLookup string = "cookie:" + CookieJWT
const CookieAuthScheme string = "="

func Protected(signingKey []byte, tokenLookup string) fiber.Handler {
	var authScheme string = HeaderAuthScheme
	if tokenLookup == "" {
		tokenLookup = HeaderTokenLookup
	}

	if tokenLookup == CookieTokenLookup {
		authScheme = CookieAuthScheme
	}
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: signingKey},
		Claims:       &service.UserClaims{},
		TokenLookup:  tokenLookup,
		AuthScheme:   authScheme,
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
