package http

import (
	"time"

	"findex/backend-core/internal/domain"
	"findex/backend-core/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	Auth       *AuthHandler
	Master     *MasterHandler
	Schedule   *ScheduleHandler
	Report     *ReportHandler
	Attendance *AttendanceHandler
}

func RegisterRoutes(app *fiber.App, h *Handlers, db *pgxpool.Pool, rdb *redis.Client, aiURL string, jwtSecret string) {
	// Health check (Public)
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "down"
		if db != nil && db.Ping(c.Context()) == nil {
			dbStatus = "up"
		}

		redisStatus := "down"
		if rdb != nil && rdb.Ping(c.Context()).Err() == nil {
			redisStatus = "up"
		}

		return c.JSON(fiber.Map{
			"app":        "Findex Backend Core",
			"status":     "healthy",
			"timestamp":  time.Now().Format(time.RFC3339),
			"database":   dbStatus,
			"redis":      redisStatus,
			"ai_service": aiURL,
		})
	})

	api := app.Group("/api/v1")

	// =========================================================================
	// PUBLIC ROUTES
	// =========================================================================
	api.Post("/auth/login", h.Auth.Login)

	// Backward compatibility alias for public/dev branch listing
	api.Get("/branches", h.Master.ListLocations)
	api.Post("/branches", h.Master.CreateLocation)

	// Attendance check-in / check-out (Can be called by employee with or without token in dev)
	api.Post("/attendance/check-in", h.Attendance.CheckIn)
	api.Post("/attendance/check-out", h.Attendance.CheckOut)
	api.Post("/attendance/self-enroll-face", h.Attendance.SelfEnrollFace)
	api.Get("/attendance/history", h.Attendance.History)

	// =========================================================================
	// AUTHENTICATED & RBAC PROTECTED MODULES
	// =========================================================================
	authRequired := api.Group("", middleware.AuthJWT(jwtSecret), middleware.TenantRLS(db))

	// -------------------------------------------------------------------------
	// 1. KELOMPOK DATA MASTER
	// -------------------------------------------------------------------------
	master := authRequired.Group("/master")

	// 1.A. Data Company -> OWNER ONLY (Bisa membuka seluruh company)
	companies := master.Group("/companies", middleware.RequireRoles(domain.RoleOwner))
	companies.Get("/", h.Master.ListCompanies)
	companies.Post("/", h.Master.CreateCompany)
	companies.Get("/:id", h.Master.GetCompany)
	companies.Put("/:id", h.Master.UpdateCompany)

	// 1.B. Data Location / Lokasi Absen -> OWNER & ADMIN COMPANY
	locations := master.Group("/locations", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdminCompany))
	locations.Get("/", h.Master.ListLocations)
	locations.Post("/", h.Master.CreateLocation)
	locations.Put("/:id", h.Master.UpdateLocation)
	locations.Delete("/:id", h.Master.DeleteLocation)

	// 1.C. Data Employee & Biometric Face -> OWNER & ADMIN COMPANY
	employees := master.Group("/employees", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdminCompany))
	employees.Get("/", h.Master.ListEmployees)
	employees.Post("/", h.Master.CreateEmployee)
	employees.Get("/:id", h.Master.GetEmployee)
	employees.Put("/:id", h.Master.UpdateEmployee)
	employees.Delete("/:id", h.Master.DeleteEmployee)
	employees.Post("/:id/face-biometric", h.Master.EnrollFaceBiometric)

	// -------------------------------------------------------------------------
	// 2. MANAJEMEN JADWAL (Shift, Batas Check-in & Awal Checkout)
	// -> OWNER & ADMIN COMPANY
	// -------------------------------------------------------------------------
	schedules := authRequired.Group("/schedules", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdminCompany))
	schedules.Get("/shifts", h.Schedule.ListShifts)
	schedules.Post("/shifts", h.Schedule.CreateShift)
	schedules.Put("/shifts/:id", h.Schedule.UpdateShift)
	schedules.Delete("/shifts/:id", h.Schedule.DeleteShift)
	schedules.Post("/assign", h.Schedule.AssignShift)

	// -------------------------------------------------------------------------
	// 3. KELOMPOK LAPORAN
	// -> OWNER & ADMIN COMPANY
	// -------------------------------------------------------------------------
	reports := authRequired.Group("/reports", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdminCompany))
	reports.Get("/attendance", h.Report.GetAttendanceReport)
}
