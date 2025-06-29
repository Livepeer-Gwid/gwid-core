// Package controllers handle incoming requests
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/types"
)

type AuthController interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService: authService}
}

func (s *authController) SignUp(c *gin.Context) {
	signupReq := c.MustGet("validatedInput").(types.SignupReq)

	user := &models.User{
		Name:     signupReq.Name,
		Email:    signupReq.Email,
		Password: signupReq.Password,
		Role:     models.Regular,
	}

	authRes, err := s.authService.SignUp(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    authRes,
	})
}

func (s *authController) Login(c *gin.Context) {
	loginReq := c.MustGet("validatedInput").(types.LoginReq)

	authRes, err := s.authService.Login(loginReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    authRes,
	})
}
