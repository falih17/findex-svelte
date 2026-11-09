package middleware

import (
	"strings"

	"findex/backend-core/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthJWT(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Error(c, fiber.StatusUnauthorized, "Token otorisasi diperlukan")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, "Token tidak valid atau telah kedaluwarsa")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Klaim token tidak valid")
		}

		c.Locals("employee_id", claims["employee_id"])
		c.Locals("company_id", claims["company_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}
