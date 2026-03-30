package service

import (
	"fmt"
	"go-back/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	secretKey string
	tokenTTL  time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, secretKey string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		secretKey: secretKey,
		tokenTTL:  24 * time.Hour,
	}
}

func (s *AuthService) AuthenticateUser(login, password string) (*repository.User, error) {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

