package database

import (
	"github.com/RomanGhost/buratino_bot.git/internal/account/database/model"
	"gorm.io/gorm"
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
		{RoleName: "admin"},
		{RoleName: "user"},
		{RoleName: "moderator"},
	}

	for _, role := range userRoles {
		if err := db.FirstOrCreate(&role, model.UserRole{RoleName: role.RoleName}).Error; err != nil {
			return err
		}
	}

	goods := []model.GoodsPrice{
		{Name: "1m vpn", Price: 20},          // 864 rub/month	0.02 rub
		{Name: "1h vpn", Price: 1000},        // 720 rub/month	1 rub
		{Name: "1d vpn", Price: 20000},       // 600 rub/month	20 rub
		{Name: "1month vpn", Price: 560000},  // 560 rub/month	560 rub
		{Name: "3month vpn", Price: 1600000}, // 533 rub/month	1600 rub
	}

	for _, good := range goods {
		if err := db.FirstOrCreate(&good, good).Error; err != nil {
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
