// Package controllers handle incoming requests
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/types"
)

type AuthController struct {
	authService           *services.AuthService
	userService           *services.UserService
	referralRewardService *services.ReferralRewardService
}

func NewAuthController(authService *services.AuthService, userService *services.UserService, referralRewardService *services.ReferralRewardService) *AuthController {
	return &AuthController{
		authService:           authService,
		userService:           userService,
		referralRewardService: referralRewardService,
	}
}

func (s *AuthController) SignUp(c *gin.Context) {
	signupReq := c.MustGet("validatedInput").(types.SignupReq)

	referralCode, err := s.userService.GenerateUniqueReferralCode()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "something went wrong",
		})
	}

	user := models.User{
		Name:         signupReq.Name,
		Email:        signupReq.Email,
		Password:     signupReq.Password,
		Role:         models.Regular,
		ReferralCode: referralCode,
	}

	authRes, err := s.authService.SignUp(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	if signupReq.ReferralCode != nil {
		s.referralRewardService.CreateReferralReward(*signupReq.ReferralCode)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    authRes,
	})
}

func (s *AuthController) Login(c *gin.Context) {
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

func (s *AuthController) ChangePassword(c *gin.Context) {
	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

	changePasswordReq := c.MustGet("validatedInput").(types.ChangePasswordReq)

	statusCode, err := s.authService.ChangePassword(changePasswordReq, reqUser.ID)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"message": "password changed successfully",
	})
}
