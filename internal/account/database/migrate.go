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
		{Name: "1m vpn", Price: 5},    // 216 rub/month	0.005 rub
		{Name: "1h vpn", Price: 250},  // 180 rub/month	0.25 rub
		{Name: "1d vpn", Price: 5000}, // 150 rub/month	5 rub
		// {Name: "1month vpn", Price: 140000}, // 140 rub/month	140 rub
		// {Name: "3month vpn", Price: 400000}, // 133 rub/month	400 rub
		{Name: "Пополнение", Price: -1}, // 0.001(10^-3) rub
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
