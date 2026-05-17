package services

import (
	"errors"
	"github.com/leonardo/netwatch/internal/domain"
	"github.com/leonardo/netwatch/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(user *domain.User) error {
	existing, _ := s.repo.FindByEmail(user.Email)
	if existing != nil && existing.ID != 0 {
		return errors.New("este e-mail já está cadastrado no sistema")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("falha interna ao criptografar a senha")
	}
	user.Password = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil || user.ID == 0 {
		return "", errors.New("credenciais informadas inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("credenciais informadas inválidas")
	}

	return "netwatch_pro_token_" + email, nil
}