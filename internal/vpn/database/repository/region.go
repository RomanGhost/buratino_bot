package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"gorm.io/gorm"
)

type RegionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) *RegionRepository {
	return &RegionRepository{db: db}
}

// Create creates a new region
func (r *RegionRepository) Create(region *model.Region) error {
	return r.db.Create(region).Error
}

// GetByName gets region by name
func (r *RegionRepository) GetByName(regionName string) (*model.Region, error) {
	var region model.Region
	err := r.db.Where("region_name = ?", regionName).First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetByShortName gets region by short name
func (r *RegionRepository) GetByShortName(shortName string) (*model.Region, error) {
	var region model.Region
	err := r.db.Where("short_name = ?", shortName).First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetAll gets all regions
func (r *RegionRepository) GetAll() ([]model.Region, error) {
	var regions []model.Region
	err := r.db.Find(&regions).Error
	return regions, err
}

// Update updates region
func (r *RegionRepository) Update(region *model.Region) error {
	return r.db.Save(region).Error
}

// Delete deletes region by name
func (r *RegionRepository) Delete(regionName string) error {
	return r.db.Delete(&model.Region{}, "region_name = ?", regionName).Error
}

// GetWithServers gets region with associated servers
func (r *RegionRepository) GetWithServers(regionName string) (*model.Region, error) {
	var region model.Region
	err := r.db.Preload("Servers").Where("region_name = ?", regionName).First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetAllWithServers gets all regions with associated servers
func (r *RegionRepository) GetAllWithServers() ([]model.Region, error) {
	var regions []model.Region
	err := r.db.Preload("Servers").Find(&regions).Error
	return regions, err
}

// Exists checks if region exists
func (r *RegionRepository) Exists(regionName string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Region{}).Where("region_name = ?", regionName).Count(&count).Error
	return count > 0, err
}

// Count gets total count of regions
func (r *RegionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Region{}).Count(&count).Error
	return count, err
}

// GetRegionServerCount gets count of servers in region
func (r *RegionRepository) GetRegionServerCount(regionName string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Server{}).Where("region = ?", regionName).Count(&count).Error
	return count, err
}
