package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/alfaizuna/go-sentinel/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims mendefinisikan payload data di dalam JWT token
type JWTClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims // Embed klaim standar (exp, iat, iss)
}

// JWTService adalah kontrak interface untuk urusan JWT
type JWTService interface {
	GenerateToken(email string, role string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// jwtService adalah implementasi konkret dari JWTService
type jwtService struct {
	secretKey       []byte
	expirationHours time.Duration
}

// NewJWTService adalah constructor function untuk inisialisasi JWTService
func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		secretKey:       []byte(cfg.JWTSecret),
		expirationHours: time.Duration(cfg.JWTExpirationHours) * time.Hour,
	}
}

// GenerateToken membuat token JWT baru yang ditandatangani dengan algoritma HMAC-SHA256 (HS256)
func (s *jwtService) GenerateToken(email string, role string) (string, error) {
	claims := JWTClaims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expirationHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-sentinel",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("gagal menandatangani token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken memvalidasi cryptographic signature token dan mengekstrak klaimnya
func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma signing (mencegah 'alg: none' exploit)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak valid: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Type assertion: pastikan token claims sesuai struct JWTClaims dan token masih valid
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid atau sudah kedaluwarsa")
	}

	return claims, nil
}
