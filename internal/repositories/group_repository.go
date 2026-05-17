package repositories

import (
	"github.com/leonardo/netwatch/database"
	"github.com/leonardo/netwatch/internal/domain"
)

type GroupRepository struct{}

func NewGroupRepository() *GroupRepository {
	return &GroupRepository{}
}

func (r *GroupRepository) Create(group *domain.Group) error {
	return database.DB.Create(group).Error
}

func (r *GroupRepository) GetAll() ([]domain.Group, error) {
	var groups []domain.Group
	err := database.DB.Find(&groups).Error
	return groups, err
}

func (r *GroupRepository) Delete(id uint) error {
	return database.DB.Unscoped().Delete(&domain.Group{}, id).Error
}