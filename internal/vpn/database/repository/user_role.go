package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) *UserRoleRepository {
	return &UserRoleRepository{db: db}
}

// Create creates a new user role
func (r *UserRoleRepository) Create(role *model.UserRole) error {
	return r.db.Create(role).Error
}

// GetByRoleName gets user role by role name
func (r *UserRoleRepository) GetByRoleName(roleName string) (*model.UserRole, error) {
	var role model.UserRole
	err := r.db.Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAll gets all user roles
func (r *UserRoleRepository) GetAll() ([]model.UserRole, error) {
	var roles []model.UserRole
	err := r.db.Find(&roles).Error
	return roles, err
}

// Update updates user role
func (r *UserRoleRepository) Update(role *model.UserRole) error {
	return r.db.Save(role).Error
}

// Delete deletes user role by role name
func (r *UserRoleRepository) Delete(roleName string) error {
	return r.db.Delete(&model.UserRole{}, "role_name = ?", roleName).Error
}

// GetWithUsers gets user role with associated users
func (r *UserRoleRepository) GetWithUsers(roleName string) (*model.UserRole, error) {
	var role model.UserRole
	err := r.db.Preload("Users").Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Exists checks if user role exists
func (r *UserRoleRepository) Exists(roleName string) (bool, error) {
	var count int64
	err := r.db.Model(&model.UserRole{}).Where("role_name = ?", roleName).Count(&count).Error
	return count > 0, err
}
