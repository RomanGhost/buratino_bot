package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) Create(role *model.UserRole) error {
	return r.db.Create(role).Error
}

func (r *UserRoleRepository) FindByName(name string) (*model.UserRole, error) {
	var role model.UserRole
	err := r.db.First(&role, "role_name = ?", name).Error
	return &role, err
}

func (r *UserRoleRepository) FindAll() ([]model.UserRole, error) {
	var roles []model.UserRole
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *UserRoleRepository) Update(role *model.UserRole) error {
	return r.db.Save(role).Error
}

func (r *UserRoleRepository) Delete(name string) error {
	return r.db.Delete(&model.UserRole{}, "role_name = ?", name).Error
}

// Загрузка роли вместе с пользователями
func (r *UserRoleRepository) FindWithUsers(name string) (*model.UserRole, error) {
	var role model.UserRole
	err := r.db.Preload("Users").First(&role, "role_name = ?", name).Error
	return &role, err
}
