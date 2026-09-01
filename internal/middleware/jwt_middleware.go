package middleware

import (
	"net/http"
	"strings"

	"github.com/alfaizuna/go-sentinel/internal/service"
	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware memvalidasi token JWT pada header 'Authorization: Bearer <token>'
func JWTMiddleware(jwtService service.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak: Header Authorization tidak ditemukan",
			})
		}

		// 2. Validasi format harus diawali dengan 'Bearer '
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Format Authorization harus 'Bearer <token>'",
			})
		}

		tokenString := parts[1]

		// 3. Validasi keabsahan token & signature dengan JWT Service
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// 4. Simpan claims ke c.Locals (ekuivalen SecurityContextHolder di Spring)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		// 5. Lanjutkan request ke handler berikutnya
		return c.Next()
	}
}

// RequireRole membatasi akses endpoint hanya untuk role tertentu (RBAC)
func RequireRole(allowedRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || role != allowedRole {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak: Anda tidak memiliki izin untuk mengakses resource ini",
			})
		}
		return c.Next()
	}
}
