package database

import (
	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"gorm.io/gorm"
)

// AutoMigrate creates all tables
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.UserRole{},
		&model.Region{},
		&model.User{},
		&model.Server{},
		&model.Key{},
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

	// Seed Regions
	regions := []model.Region{
		{RegionName: "Netherlands", ShortName: "NL"},
		{RegionName: "Moscow", ShortName: "RU"},
		{RegionName: "Germany", ShortName: "DE"},
	}

	for _, region := range regions {
		if err := db.FirstOrCreate(&region, model.Region{RegionName: region.RegionName}).Error; err != nil {
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
