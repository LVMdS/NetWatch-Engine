package repositories

import (
	"github.com/leonardo/netwatch/database"
	"github.com/leonardo/netwatch/internal/domain"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *domain.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := database.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) Delete(id uint) error {
	// Hard Delete garante que o e-mail seja removido de vez do arquivo netwatch.db
	return database.DB.Unscoped().Delete(&domain.User{}, id).Error
}