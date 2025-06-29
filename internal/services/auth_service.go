// Package services declars bussiness logic
package services

import (
	"errors"

	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/types"
)

type AuthService interface {
	SignUp(user *models.User) (types.AuthRes, error)
	Login(loginReq types.LoginReq) (types.AuthRes, error)
}

type authService struct {
	userRepository repositories.UserRepository
	jwtService     JWTService
}

func NewAuthService(userRepository repositories.UserRepository, jwtService JWTService) AuthService {
	return &authService{userRepository: userRepository, jwtService: jwtService}
}

func (s *authService) SignUp(user *models.User) (types.AuthRes, error) {
	_, result := s.userRepository.FindByEmail(user.Email)

	if result.RowsAffected > 0 {
		return types.AuthRes{}, errors.New("user already exists")
	}

	err := s.userRepository.CreateUser(user)
	if err != nil {
		return types.AuthRes{}, err
	}

	tokenString, err := s.jwtService.SignJWT(user)

	return types.AuthRes{
		ID:          user.ID,
		Role:        string(user.Role),
		AccessToken: tokenString,
	}, err
}

func (s *authService) Login(loginReq types.LoginReq) (types.AuthRes, error) {
	user, result := s.userRepository.FindByEmail(loginReq.Email)

	if result.RowsAffected == 0 {
		return types.AuthRes{}, errors.New("invalid credentials")
	}

	if err := user.CheckPassword(loginReq.Password); err != nil {
		return types.AuthRes{}, errors.New("invalid credentials")
	}

	tokenString, err := s.jwtService.SignJWT(user)

	return types.AuthRes{
		ID:          user.ID,
		Role:        string(user.Role),
		AccessToken: tokenString,
	}, err
}
