package handler

import (
	"net/http"

	"github.com/alfaizuna/go-sentinel/internal/dto"
	"github.com/alfaizuna/go-sentinel/internal/service"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler menangani HTTP request untuk endpoint autentikasi
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler adalah constructor function untuk AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register menangani POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// 1. Parse JSON body ke struct RegisterRequest (mirip @RequestBody di Spring)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request body tidak valid",
		})
	}

	// 2. Validasi input sederhana
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, dan password wajib diisi",
		})
	}

	// 3. Panggil service layer
	res, err := h.authService.Register(req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"data":    res,
	})
}

// Login menangani POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// 1. Parse JSON body ke struct LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request body tidak valid",
		})
	}

	// 2. Validasi input
	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email dan password wajib diisi",
		})
	}

	// 3. Panggil service layer
	res, err := h.authService.Login(req)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"data":    res,
	})
}
