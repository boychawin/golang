package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func ValidateJWT(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Token")
		var SECRET = []byte(viper.GetString("jwt_ACCESS_TOKEN_SECRET"))

		if authHeader != "" {
			// Parse the JWT token
			parsedToken, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
				// Verify the signing method and key
				if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, fmt.Errorf("invalid signing method")
				}
				return SECRET, nil
			})

			if err != nil {
				return c.JSON(fiber.Map{
					"Expired":  true,
					"Status":   "error",
					"Messages": err.Error(),
				})
			}

			if parsedToken.Valid {
				// Access the token claims (payload)
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					return c.JSON(fiber.Map{
						"Expired":  true,
						"Status":   "error",
						"Messages": err.Error(),
					})
				}
				// Store the claims as context data
				c.Locals("claims", claims)
				// Call the next middleware or handler with the payload
				return next(c)
			}
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Expired":  false,
			"Status":   "error",
			"Messages": "Token is required",
		})
	}
}