package database

import (
	"fmt"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AutoMigrate creates all tables
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Provider{},
		&model.Region{},
		&model.User{},
		&model.Server{},
		&model.Key{},
	)
}

// SeedData populates initial data
func SeedData(db *gorm.DB) error {
	// Seed Regions
	regions := []model.Region{
		{RegionName: "Netherlands", ShortName: "NL"},
		{RegionName: "Moscow", ShortName: "RU"},
		{RegionName: "Germany", ShortName: "DE"},
	}

	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&regions).Error
	if err != nil {
		return fmt.Errorf("failed to initialize actions: %v", err)
	}

	providers := []model.Provider{
		model.Outline, model.Wireguard,
	}

	err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&providers).Error
	if err != nil {
		return fmt.Errorf("failed to initialize actions: %v", err)
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
