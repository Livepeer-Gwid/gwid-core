// Package di mounts all services
package di

import (
	"gorm.io/gorm"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/controllers"
	"gwid.io/gwid-core/internal/database"
	"gwid.io/gwid-core/internal/repositories"
	"gwid.io/gwid-core/internal/services"
)

type Container struct {
	DB             *gorm.DB
	UserRepository repositories.UserRepository
	AuthService    services.AuthService
	JwtService     services.JWTService
	AuthController controllers.AuthController
}

func NewContainer(conf *config.Config) *Container {
	database.InitDB(conf)

	db := database.DB

	jwtService := services.NewJwtService(conf)

	userRepository := repositories.NewUserRepository(db)

	authService := services.NewAuthService(userRepository, jwtService)

	authController := controllers.NewAuthController(authService)

	return &Container{
		DB:             db,
		UserRepository: userRepository,
		AuthService:    authService,
		JwtService:     jwtService,
		AuthController: authController,
	}
}
