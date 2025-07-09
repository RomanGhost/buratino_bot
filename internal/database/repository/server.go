package repository

import (
	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

// Create creates a new server
func (r *ServerRepository) Create(server *model.Server) error {
	return r.db.Create(server).Error
}

// GetByID gets server by ID
func (r *ServerRepository) GetByID(id uint) (*model.Server, error) {
	var server model.Server
	err := r.db.First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// GetByRegion gets servers by region
func (r *ServerRepository) GetByRegion(region string) ([]model.Server, error) {
	var servers []model.Server
	err := r.db.Where("region = ?", region).Find(&servers).Error
	return servers, err
}

// GetAll gets all servers with pagination
func (r *ServerRepository) GetAll(offset, limit int) ([]model.Server, error) {
	var servers []model.Server
	err := r.db.Offset(offset).Limit(limit).Find(&servers).Error
	return servers, err
}

// GetByIPv4 gets server by IPv4 address
func (r *ServerRepository) GetByIPv4(ipv4 string) (*model.Server, error) {
	var server model.Server
	err := r.db.Where("ipv4 = ?", ipv4).First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// GetByIPv6 gets server by IPv6 address
func (r *ServerRepository) GetByIPv6(ipv6 string) (*model.Server, error) {
	var server model.Server
	err := r.db.Where("ipv6 = ?", ipv6).First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// Update updates server
func (r *ServerRepository) Update(server *model.Server) error {
	return r.db.Save(server).Error
}

// Delete deletes server (soft delete)
func (r *ServerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Server{}, id).Error
}

// GetWithKeys gets server with associated keys
func (r *ServerRepository) GetWithKeys(id uint) (*model.Server, error) {
	var server model.Server
	err := r.db.Preload("Keys").First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// GetWithRegion gets server with region info
func (r *ServerRepository) GetWithRegion(id uint) (*model.Server, error) {
	var server model.Server
	err := r.db.Preload("RegionInfo").First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// GetWithFullInfo gets server with all associations
func (r *ServerRepository) GetWithFullInfo(id uint) (*model.Server, error) {
	var server model.Server
	err := r.db.Preload("RegionInfo").Preload("Keys").First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

// Count gets total count of servers
func (r *ServerRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Server{}).Count(&count).Error
	return count, err
}

// CountByRegion gets count of servers by region
func (r *ServerRepository) CountByRegion(region string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Server{}).Where("region = ?", region).Count(&count).Error
	return count, err
}

// GetAvailableServers gets servers with available capacity (example logic)
func (r *ServerRepository) GetAvailableServers(region string) ([]model.Server, error) {
	var servers []model.Server
	// Пример логики - серверы с менее чем 100 активными ключами
	err := r.db.Where("region = ?", region).
		Joins("LEFT JOIN keys ON servers.id = keys.server_id AND keys.is_active = true").
		Group("servers.id").
		Having("COUNT(keys.id) < ?", 100).
		Find(&servers).Error
	return servers, err
}

// GetServerLoad gets server load (count of active keys)
func (r *ServerRepository) GetServerLoad(id uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Key{}).Where("server_id = ? AND is_active = true", id).Count(&count).Error
	return count, err
}

// ExistsByIPv4 checks if server exists by IPv4
func (r *ServerRepository) ExistsByIPv4(ipv4 string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Server{}).Where("ipv4 = ?", ipv4).Count(&count).Error
	return count > 0, err
}

// ExistsByIPv6 checks if server exists by IPv6
func (r *ServerRepository) ExistsByIPv6(ipv6 string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Server{}).Where("ipv6 = ?", ipv6).Count(&count).Error
	return count > 0, err
}
