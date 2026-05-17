package repositories

import (
	"time"
	"github.com/leonardo/netwatch/database"
	"github.com/leonardo/netwatch/internal/domain"
	"gorm.io/gorm"
)

type DeviceRepository struct{}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{}
}

func (r *DeviceRepository) Create(device *domain.Device) error {
	return database.DB.Create(device).Error
}

func (r *DeviceRepository) GetAll() ([]domain.Device, error) {
	var devices []domain.Device
	err := database.DB.Preload("Group").Find(&devices).Error
	return devices, err
}

func (r *DeviceRepository) Update(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&domain.Device{}).Where("id = ?", id).Updates(updates).Error
}

func (r *DeviceRepository) Delete(id uint) error {
	return database.DB.Unscoped().Delete(&domain.Device{}, id).Error
}

// RecordCheck consolida o estado atual e alimenta a métrica de SLA
func (r *DeviceRepository) RecordCheck(ip string, checkPassed bool, consolidatedStatus string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"total_checks": gorm.Expr("total_checks + 1"),
		"status":       consolidatedStatus,
	}
	if checkPassed {
		updates["successful_checks"] = gorm.Expr("successful_checks + 1")
		updates["last_seen"] = &now
	}
	return database.DB.Model(&domain.Device{}).Where("ip_address = ?", ip).Updates(updates).Error
}