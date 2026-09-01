package dto

// RegisterRequest adalah struktur data untuk request registrasi
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest adalah struktur data untuk request login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse adalah struktur data untuk response token autentikasi
type AuthResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
