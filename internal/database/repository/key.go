package repository

import (
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/database/model"
	"gorm.io/gorm"
)

type KeyRepository struct {
	db *gorm.DB
}

func NewKeyRepository(db *gorm.DB) *KeyRepository {
	return &KeyRepository{db: db}
}

// Create creates a new key
func (r *KeyRepository) Create(key *model.Key) error {
	return r.db.Create(key).Error
}

// GetByID gets key by ID
func (r *KeyRepository) GetByID(id uint) (*model.Key, error) {
	var key model.Key
	err := r.db.First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetByUserID gets keys by user ID
func (r *KeyRepository) GetByUserID(userID uint) ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Where("user_id = ?", userID).Find(&keys).Error
	return keys, err
}

// GetByServerID gets keys by server ID
func (r *KeyRepository) GetByServerID(serverID uint) ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Where("server_id = ?", serverID).Find(&keys).Error
	return keys, err
}

// GetActiveKeys gets all active keys
func (r *KeyRepository) GetActiveKeys() ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Where("is_active = ?", true).Find(&keys).Error
	return keys, err
}

// GetActiveKeysByUser gets active keys by user ID
func (r *KeyRepository) GetActiveKeysByUser(userID uint) ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&keys).Error
	return keys, err
}

// GetActiveKeysByServer gets active keys by server ID
func (r *KeyRepository) GetActiveKeysByServer(serverID uint) ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Where("server_id = ? AND is_active = ?", serverID, true).Find(&keys).Error
	return keys, err
}

// GetExpiredKeys gets expired keys
func (r *KeyRepository) GetExpiredKeys() ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Where("deadline_time < ?", time.Now()).Find(&keys).Error
	return keys, err
}

// GetExpiringSoon gets keys expiring within specified duration
func (r *KeyRepository) GetExpiringSoon(duration time.Duration) ([]model.Key, error) {
	var keys []model.Key
	expiryTime := time.Now().Add(duration)
	err := r.db.Where("deadline_time BETWEEN ? AND ? AND is_active = ?",
		time.Now(), expiryTime, true).Find(&keys).Error
	return keys, err
}

// GetAll gets all keys with pagination
func (r *KeyRepository) GetAll(offset, limit int) ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Offset(offset).Limit(limit).Find(&keys).Error
	return keys, err
}

// Update updates key
func (r *KeyRepository) Update(key *model.Key) error {
	return r.db.Save(key).Error
}

// Delete deletes key (soft delete)
func (r *KeyRepository) Delete(id uint) error {
	return r.db.Delete(&model.Key{}, id).Error
}

// DeactivateKey deactivates key
func (r *KeyRepository) DeactivateKey(id uint) error {
	return r.db.Model(&model.Key{}).Where("id = ?", id).Update("is_active", false).Error
}

// ActivateKey activates key
func (r *KeyRepository) ActivateKey(id uint) error {
	return r.db.Model(&model.Key{}).Where("id = ?", id).Update("is_active", true).Error
}

// ExtendKey extends key deadline
func (r *KeyRepository) ExtendKey(id uint, newDeadline time.Time) error {
	return r.db.Model(&model.Key{}).Where("id = ?", id).Update("deadline_time", newDeadline).Error
}

// GetWithUser gets key with user info
func (r *KeyRepository) GetWithUser(id uint) (*model.Key, error) {
	var key model.Key
	err := r.db.Preload("User").First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetWithServer gets key with server info
func (r *KeyRepository) GetWithServer(id uint) (*model.Key, error) {
	var key model.Key
	err := r.db.Preload("Server").First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetWithFullInfo gets key with all associations
func (r *KeyRepository) GetWithFullInfo(id uint) (*model.Key, error) {
	var key model.Key
	err := r.db.Preload("User").Preload("Server").First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// Count gets total count of keys
func (r *KeyRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Key{}).Count(&count).Error
	return count, err
}

// CountActive gets count of active keys
func (r *KeyRepository) CountActive() (int64, error) {
	var count int64
	err := r.db.Model(&model.Key{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

// CountByUser gets count of keys by user ID
func (r *KeyRepository) CountByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Key{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountActiveByUser gets count of active keys by user ID
func (r *KeyRepository) CountActiveByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Key{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&count).Error
	return count, err
}

// CountByServer gets count of keys by server ID
func (r *KeyRepository) CountByServer(serverID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Key{}).Where("server_id = ?", serverID).Count(&count).Error
	return count, err
}

// CountActiveByServer gets count of active keys by server ID
func (r *KeyRepository) CountActiveByServer(serverID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Key{}).Where("server_id = ? AND is_active = ?", serverID, true).Count(&count).Error
	return count, err
}

// CleanupExpiredKeys deactivates expired keys
func (r *KeyRepository) CleanupExpiredKeys() error {
	return r.db.Model(&model.Key{}).
		Where("deadline_time < ? AND is_active = ?", time.Now(), true).
		Update("is_active", false).Error
}

// GetUserKeysByRegion gets user's keys by region
func (r *KeyRepository) GetUserKeysByRegion(userID uint, region string) ([]model.Key, error) {
	var keys []model.Key
	err := r.db.Joins("JOIN servers ON keys.server_id = servers.id").
		Where("keys.user_id = ? AND servers.region = ?", userID, region).
		Find(&keys).Error
	return keys, err
}
