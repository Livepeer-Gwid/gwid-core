package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/types"
)

type JwtService struct {
	config *config.Config
}

func NewJwtService(config *config.Config) *JwtService {
	return &JwtService{
		config: config,
	}
}

func (s *JwtService) SignJWT(user *models.User) (string, error) {
	claims := &types.JwtCustomClaims{
		ID:   user.ID,
		Role: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gwid-core",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JwtSecret))

	return tokenString, err
}
