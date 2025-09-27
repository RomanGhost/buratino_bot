package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Wallet").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByTelegramID(tgID int64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Wallet").First(&user, "telegram_id = ?", tgID).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
