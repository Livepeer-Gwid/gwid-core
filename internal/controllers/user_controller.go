package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/services"
)

type UserController interface {
	GetCurrentUserProfile(c *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

func (s *userController) GetCurrentUserProfile(c *gin.Context) {
	reqUser := c.MustGet("user").(*services.JwtCustomClaims)

	user, err := s.userService.GetCurrentUserProfile(reqUser.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}
