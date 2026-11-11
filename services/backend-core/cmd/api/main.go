package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"findex/backend-core/internal/config"
	delivery "findex/backend-core/internal/delivery/http"
	"findex/backend-core/internal/pkg/aiservice"
	"findex/backend-core/internal/pkg/database"
	"findex/backend-core/internal/repository"
	"findex/backend-core/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// 1. Initialize PostgreSQL Connection Pool
	dbPool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("Warning: Failed to establish DB pool at startup: %v", err)
	} else {
		log.Println("✓ Connected to PostgreSQL with PostGIS & pgvector successfully")
	}

	// 2. Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,
	})

	// 3. Initialize AI Biometric Service Client
	aiClient := aiservice.NewClient(cfg.AIServiceURL)

	// 4. Initialize Repositories
	companyRepo := repository.NewCompanyRepository(dbPool)
	branchRepo := repository.NewBranchRepository(dbPool)
	empRepo := repository.NewEmployeeRepository(dbPool)
	shiftRepo := repository.NewShiftRepository(dbPool)
	reportRepo := repository.NewReportRepository(dbPool)
	attRepo := repository.NewAttendanceRepository(dbPool)

	// 5. Initialize Usecases
	authUsecase := usecase.NewAuthUsecase(empRepo, companyRepo, cfg.JWTSecret)
	masterUsecase := usecase.NewMasterUsecase(companyRepo, branchRepo, empRepo, aiClient)
	scheduleUsecase := usecase.NewScheduleUsecase(shiftRepo)
	reportUsecase := usecase.NewReportUsecase(reportRepo)
	attUsecase := usecase.NewAttendanceUsecase(attRepo, branchRepo, empRepo, shiftRepo, companyRepo, aiClient)

	// 6. Initialize Handlers
	handlers := &delivery.Handlers{
		Auth:       delivery.NewAuthHandler(authUsecase),
		Master:     delivery.NewMasterHandler(masterUsecase),
		Schedule:   delivery.NewScheduleHandler(scheduleUsecase),
		Report:     delivery.NewReportHandler(reportUsecase),
		Attendance: delivery.NewAttendanceHandler(attUsecase),
	}

	// 7. Setup Fiber App
	app := fiber.New(fiber.Config{
		AppName:      "Findex Multi-Tenant SaaS Backend Core",
		BodyLimit:    15 * 1024 * 1024, // 15MB limit for photo uploads
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Company-ID",
	}))

	// Register Routes with RBAC & Middlewares
	delivery.RegisterRoutes(app, handlers, dbPool, rdb, cfg.AIServiceURL, cfg.JWTSecret)

	log.Printf("🚀 Findex Backend Core running on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
