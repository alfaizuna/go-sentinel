package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB membuat dan mengembalikan koneksi instance *gorm.DB
func ConnectDB(cfg *Config) (*gorm.DB, error) {
	// Format Data Source Name (DSN) untuk PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	// Membuka koneksi database dengan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal terkoneksi ke database: %w", err)
	}

	log.Println("✅ Berhasil terkoneksi ke database PostgreSQL!")
	return db, nil
}
