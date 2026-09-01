package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menampung semua variabel konfigurasi aplikasi
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret          string
	JWTExpirationHours int

	// App
	AppPort string
}

// LoadConfig membaca file .env dan mengembalikannya sebagai pointer struct *Config
func LoadConfig() *Config {
	// Load file .env jika ada (jika di production/docker, env dibaca langsung dari sistem OS)
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan, membaca dari OS environment variables...")
	}

	// Helper konversi string ke int untuk JWT Expiration
	jwtHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		jwtHours = 24
	}

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "gosentinel_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "default_secret_key_change_me_in_production"),
		JWTExpirationHours: jwtHours,
		AppPort:            getEnv("APP_PORT", "8081"),
	}
}

// getEnv adalah helper function sederhana dengan fallback default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
