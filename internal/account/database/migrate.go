package database

import (
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AutoMigrate creates all tables
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.UserRole{},
		&model.GoodsPrice{},
		&model.Operation{},
		&model.Wallet{},
	)
}

// SeedData populates initial data
func SeedData(db *gorm.DB) error {
	// Seed model.UserRoles
	userRoles := []model.UserRole{
		model.CommonUserRole,
		model.AdminRole,
		model.ModeratorRole,
	}

	for _, role := range userRoles {
		if err := db.FirstOrCreate(&role, model.UserRole{RoleName: role.RoleName}).Error; err != nil {
			return err
		}
	}

	goods := []model.GoodsPrice{
		model.VPN1Min,
		model.VPN1Hour,
		model.VPN1Day,
		model.VPN1Month,
		model.TopUP,
	}

	for _, good := range goods {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sys_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "price"}),
		}).Create(&good).Error; err != nil {
			return err
		}
	}

	return nil
}

// InitDB initializes database with migrations and seed data
func InitDB(db *gorm.DB) error {
	// Run migrations
	if err := AutoMigrate(db); err != nil {
		return err
	}

	// Seed initial data
	if err := SeedData(db); err != nil {
		return err
	}

	return nil
}
