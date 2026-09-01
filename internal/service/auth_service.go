package service

import (
	"errors"
	"fmt"

	"github.com/alfaizuna/go-sentinel/internal/dto"
	"github.com/alfaizuna/go-sentinel/internal/model"
	"github.com/alfaizuna/go-sentinel/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService adalah kontrak interface untuk operasi autentikasi
type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtService JWTService
}

// NewAuthService adalah constructor function dengan Dependency Injection manual
func NewAuthService(userRepo repository.UserRepository, jwtService JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register mendaftarkan user baru, melakukan hashing password, dan mengembalikan token
func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 1. Validasi apakah email sudah terdaftar
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, errors.New("email sudah digunakan")
	}

	// 2. Hash password dengan Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	// 3. Buat entity user baru
	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     model.RoleUser, // Default role: USER
	}

	// 4. Simpan ke database
	if err := s.userRepo.Save(&newUser); err != nil {
		return nil, fmt.Errorf("gagal menyimpan user: %w", err)
	}

	// 5. Generate JWT token
	token, err := s.jwtService.GenerateToken(newUser.Email, string(newUser.Role))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		Email: newUser.Email,
		Role:  string(newUser.Role),
	}, nil
}

// Login memvalidasi kredensial user dan mengembalikan token jika valid
func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("terjadi kesalahan sistem: %w", err)
	}
	if user == nil {
		return nil, errors.New("email atau password salah")
	}

	// 2. Komparasi password mentah dengan password hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau password salah")
	}

	// 3. Generate JWT token
	token, err := s.jwtService.GenerateToken(user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		Email: user.Email,
		Role:  string(user.Role),
	}, nil
}
