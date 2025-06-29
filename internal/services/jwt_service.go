package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/models"
)

type JWTService interface {
	SignJWT(user *models.User) (string, error)
}

type jwtService struct {
	configService *config.Config
}

func NewJwtService(configService *config.Config) JWTService {
	return &jwtService{configService: configService}
}

type JwtCustomClaims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
	jwt.RegisteredClaims
}

func (s *jwtService) SignJWT(user *models.User) (string, error) {
	claims := &JwtCustomClaims{
		ID:   user.ID,
		Role: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lolarpay",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.configService.JwtSecret))

	return tokenString, err
}
