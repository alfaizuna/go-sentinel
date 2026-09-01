package main

import (
	"fmt"
	"log"

	"github.com/alfaizuna/go-sentinel/internal/config"
	"github.com/alfaizuna/go-sentinel/internal/handler"
	"github.com/alfaizuna/go-sentinel/internal/middleware"
	"github.com/alfaizuna/go-sentinel/internal/model"
	"github.com/alfaizuna/go-sentinel/internal/repository"
	"github.com/alfaizuna/go-sentinel/internal/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Database Connection & Migration
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// 3. Wiring Dependency Injection
	userRepo := repository.NewUserRepository(db)
	jwtService := service.NewJWTService(cfg)
	authService := service.NewAuthService(userRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)
	demoHandler := handler.NewDemoHandler()

	// 4. Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		AppName: "Go Sentinel - Auth & Security API",
	})

	// 5. Setup Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "UP",
			"message": "Go Sentinel is running smoothly",
		})
	})

	// API Group: /api/v1
	api := app.Group("/api/v1")

	// Public Routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected Routes (Diproteksi oleh JWTMiddleware)
	demo := api.Group("/demo", middleware.JWTMiddleware(jwtService))
	demo.Get("/hello", demoHandler.Hello)
	demo.Get("/admin", middleware.RequireRole("ADMIN"), demoHandler.AdminOnly)

	// 6. Start Server
	listenAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on port %s", cfg.AppPort)

	if err := app.Listen(listenAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
