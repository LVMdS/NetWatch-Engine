package services

import (
	"github.com/leonardo/netwatch/internal/domain"
	"github.com/leonardo/netwatch/internal/repositories"
)

type DeviceService struct {
	repo *repositories.DeviceRepository
}

func NewDeviceService(repo *repositories.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) ListAll() ([]domain.Device, error) {
	return s.repo.GetAll()
}

func (s *DeviceService) RegisterDevice(ip, name string) error {
	device := &domain.Device{
		IPAddress: ip,
		Name:      name,
		Status:    "UNKNOWN",
	}
	return s.repo.Create(device)
}