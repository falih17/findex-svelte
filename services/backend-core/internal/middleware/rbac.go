package middleware

import (
	"fmt"

	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// RequireRoles verifies whether the authenticated user has one of the allowed roles.
// Owner always has access to all modules.
func RequireRoles(allowedRoles ...domain.EmployeeRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("role")
		if roleVal == nil {
			return response.Error(c, fiber.StatusUnauthorized, "Sesi autentikasi tidak valid atau belum login")
		}

		userRole := domain.EmployeeRole(fmt.Sprintf("%v", roleVal))

		// OWNER has master access to all modules
		if userRole == domain.RoleOwner || userRole == domain.RoleSuperAdmin {
			return c.Next()
		}

		// Check if user's role is in allowed list
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				return c.Next()
			}
			// Backward compatibility: map HR_ADMIN to ADMIN_COMPANY
			if allowed == domain.RoleAdminCompany && userRole == domain.RoleHRAdmin {
				return c.Next()
			}
		}

		return response.Error(c, fiber.StatusForbidden, 
			fmt.Sprintf("Akses ditolak: role '%s' tidak memiliki izin untuk modul ini", userRole))
	}
}
