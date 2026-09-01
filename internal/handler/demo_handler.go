package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type DemoHandler struct{}

func NewDemoHandler() *DemoHandler {
	return &DemoHandler{}
}

// Hello hanya bisa diakses jika user sudah login (membawa token valid)
func (h *DemoHandler) Hello(c *fiber.Ctx) error {
	// Ambil data yang disimpan oleh JWTMiddleware sebelumnya
	userEmail := c.Locals("userEmail").(string)
	userRole := c.Locals("userRole").(string)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Halo! Kamu berhasil mengakses protected endpoint",
		"user": fiber.Map{
			"email": userEmail,
			"role":  userRole,
		},
	})
}

// AdminOnly hanya bisa diakses jika role-nya ADMIN
func (h *DemoHandler) AdminOnly(c *fiber.Ctx) error {
	userEmail := c.Locals("userEmail").(string)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Selamat datang di Halaman Rahasia Admin!",
		"admin":   userEmail,
	})
}
