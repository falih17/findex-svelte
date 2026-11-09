package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	Auth       *AuthHandler
	Branch     *BranchHandler
	Attendance *AttendanceHandler
}

func RegisterRoutes(app *fiber.App, h *Handlers, db *pgxpool.Pool, rdb *redis.Client, aiURL string) {
	// Health check
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

	// Auth
	api.Post("/auth/login", h.Auth.Login)

	// Branches (Office Locations & Geofence Settings)
	api.Get("/branches", h.Branch.List)
	api.Post("/branches", h.Branch.Create)

	// Attendance & Biometrics
	api.Post("/attendance/check-in", h.Attendance.CheckIn)
	api.Get("/attendance/history", h.Attendance.History)
}
