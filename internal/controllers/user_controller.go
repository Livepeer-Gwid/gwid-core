package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/types"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (s *UserController) GetCurrentUserProfile(c *gin.Context) {
	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

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
