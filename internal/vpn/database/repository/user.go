package repository

import (
	"errors"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID gets user by ID
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByTelegramID gets user by telegram ID
func (r *UserRepository) GetByTelegramID(telegramID int64) (*model.User, error) {
	var user model.User
	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll gets all users with pagination
func (r *UserRepository) GetAll(offset, limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// Update updates user
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete deletes user (soft delete)
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// Count gets total count of users
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

// ExistsByTelegramID checks if user exists by telegram ID
func (r *UserRepository) ExistsByTelegramID(telegramID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("telegram_id = ?", telegramID).Count(&count).Error
	return count > 0, err
}
