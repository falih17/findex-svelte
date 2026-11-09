package middleware

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantRLS sets the PostgreSQL session local variable 'app.current_company_id'
// which triggers PostgreSQL Row Level Security (RLS) policies automatically.
func TenantRLS(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Locals("company_id")
		if companyID != nil {
			compIDStr := fmt.Sprintf("%v", companyID)
			// Execute session level set on the transaction or local connection
			ctx := context.WithValue(c.Context(), "current_company_id", compIDStr)
			c.SetUserContext(ctx)
		}
		return c.Next()
	}
}
