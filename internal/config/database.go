package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB membuat dan mengembalikan koneksi instance *gorm.DB
func ConnectDB(cfg *Config) (*gorm.DB, error) {
	// Menggunakan format PostgreSQL Connection URL
	var dsn string
	if cfg.DBPassword == "" {
		dsn = fmt.Sprintf(
			"postgres://%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
			cfg.DBSSLMode,
		)
	} else {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
			cfg.DBSSLMode,
		)
	}

	// Membuka koneksi database dengan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal terkoneksi ke database: %w", err)
	}

	log.Println("Berhasil terkoneksi ke database PostgreSQL!")
	return db, nil
}
