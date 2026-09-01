package repository

import (
	"errors"

	"github.com/alfaizuna/go-sentinel/internal/model"
	"gorm.io/gorm"
)

// UserRepository adalah kontrak interface untuk akses data user
type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	ExistsByEmail(email string) bool
	Save(user *model.User) error
}

// userRepository adalah implementasi konkret dari UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository adalah constructor function untuk menginisialisasi repository (Dependency Injection manual)
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByEmail mencari user berdasarkan alamat email
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil jika user tidak ditemukan (mirip Optional.empty() di Java)
		}
		return nil, err // Return error jika ada masalah query/koneksi database
	}
	return &user, nil
}

// ExistsByEmail mengecek apakah email sudah terdaftar di database
func (r *userRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// Save menyimpan user baru atau mengupdate user yang sudah ada
func (r *userRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}
