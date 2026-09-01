package model

import (
	"gorm.io/gorm"
)

// Role merepresentasikan role user (mirip Enum di Java)
type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

// User merepresentasikan entitas tabel 'users' di database
type User struct {
	gorm.Model        // Menyediakan otomatis: ID (uint, Primary Key), CreatedAt, UpdatedAt, DeletedAt (Soft Delete)
	Name       string `gorm:"type:varchar(100);not null" json:"name"`
	Email      string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password   string `gorm:"type:varchar(255);not null" json:"-"` // json:"-" artinya field password tidak akan di-expose saat di-serialize ke JSON response
	Role       Role   `gorm:"type:varchar(20);default:'USER'" json:"role"`
}
