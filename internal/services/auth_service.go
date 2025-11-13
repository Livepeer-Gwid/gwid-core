// Package services hold business logic
package services

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/types"
)

type AuthService struct {
	userRepository *repositories.UserRepository
	jwtService     *JwtService
}

func NewAuthService(userRepository *repositories.UserRepository, jwtService *JwtService) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (s *AuthService) SignUp(user *models.User) (types.AuthRes, error) {
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

func (s *AuthService) Login(loginReq types.LoginReq) (types.AuthRes, error) {
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

func (s *AuthService) ChangePassword(changePasswordReq types.ChangePasswordReq, userID uuid.UUID) (int, error) {
	user, result := s.userRepository.FindByID(userID)

	if result.RowsAffected == 0 {
		return http.StatusNotFound, errors.New("user not found")
	}

	if err := user.CheckPassword(changePasswordReq.CurrentPassword); err != nil {
		return http.StatusBadRequest, errors.New("invalid current password")
	}

	if err := user.HashPassword(changePasswordReq.NewPassword); err != nil {
		return http.StatusInternalServerError, err
	}

	if result := s.userRepository.UpdateUser(user); result.RowsAffected == 0 {
		return http.StatusInternalServerError, errors.New("unable to change password")
	}

	return http.StatusOK, nil
}
