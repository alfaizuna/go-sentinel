package main

import (
	"fmt"
	"log"

	"github.com/alfaizuna/go-sentinel/internal/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	_ = db

	app := fiber.New(fiber.Config{
		AppName: "Go Sentinel - Auth & Security API",
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "UP",
			"message": "Go Sentinel is running smoothly 🚀",
		})
	})

	listenAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Server running on port %s", cfg.AppPort)

	if err := app.Listen(listenAddr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
