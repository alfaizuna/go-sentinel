package service_test

import (
	"testing"

	"github.com/alfaizuna/go-sentinel/internal/config"
	"github.com/alfaizuna/go-sentinel/internal/service"
)

// TestGenerateAndValidateToken_Success menguji generate token dan validasi token yang valid
func TestGenerateAndValidateToken_Success(t *testing.T) {
	// 1. Setup Mock Config
	mockCfg := &config.Config{
		JWTSecret:          "rahasia_super_aman_12345678901234567890",
		JWTExpirationHours: 24,
	}

	jwtService := service.NewJWTService(mockCfg)

	expectedEmail := "tester@example.com"
	expectedRole := "ADMIN"

	// 2. Eksekusi GenerateToken
	token, err := jwtService.GenerateToken(expectedEmail, expectedRole)
	if err != nil {
		t.Fatalf("Gagal generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Token tidak boleh kosong")
	}

	// 3. Eksekusi ValidateToken
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Token valid tapi gagal divalidasi: %v", err)
	}

	// 4. Assert Data Claims
	if claims.Email != expectedEmail {
		t.Errorf("Expected email %s, tapi didapat %s", expectedEmail, claims.Email)
	}

	if claims.Role != expectedRole {
		t.Errorf("Expected role %s, tapi didapat %s", expectedRole, claims.Role)
	}
}

// TestValidateToken_InvalidSignature menguji bahwa token dengan secret key berbeda harus ditolak
func TestValidateToken_InvalidSignature(t *testing.T) {
	// Service 1 (Secret Asli)
	service1 := service.NewJWTService(&config.Config{
		JWTSecret:          "kunci_rahasia_pertama_123456789012",
		JWTExpirationHours: 24,
	})

	// Service 2 (Secret Berbeda / Penyerang)
	service2 := service.NewJWTService(&config.Config{
		JWTSecret:          "kunci_rahasia_kedua_yang_berbeda_123",
		JWTExpirationHours: 24,
	})

	// Generate token dari service 1
	token, _ := service1.GenerateToken("attacker@example.com", "ADMIN")

	// Validasi menggunakan service 2 (harus gagal)
	_, err := service2.ValidateToken(token)
	if err == nil {
		t.Error("Harusnya validasi gagal karena secret key berbeda, tapi malah berhasil!")
	}
}

// TestValidateToken_MalformedToken menguji bahwa token dengan format rusak harus ditolak
func TestValidateToken_MalformedToken(t *testing.T) {
	jwtService := service.NewJWTService(&config.Config{
		JWTSecret:          "rahasia_12345678901234567890",
		JWTExpirationHours: 24,
	})

	invalidToken := "bukan.token.jwt.yang.benar"

	_, err := jwtService.ValidateToken(invalidToken)
	if err == nil {
		t.Error("Harusnya token rusak ditolak, tapi tidak menghasilkan error!")
	}
}
