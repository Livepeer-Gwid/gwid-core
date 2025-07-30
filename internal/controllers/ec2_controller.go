package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/types"
)

type EC2Controller struct {
	ec2Service *services.EC2Service
}

func NewEC2Controller(ec2Service *services.EC2Service) *EC2Controller {
	return &EC2Controller{
		ec2Service: ec2Service,
	}
}

func (s *EC2Controller) GetEC2InstanceTypes(c *gin.Context) {
	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

	params, exists := middleware.GetQueryParams(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get query params"})
		return
	}

	credentialsID, err := uuid.Parse(params.Filters["credentials_id"])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid credentials ID",
		})

		return
	}

	ec2, statusCode, err := s.ec2Service.GetEC2InstanceTypes(reqUser.ID, credentialsID, "eu-central-1")
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err,
		})

		return
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    ec2,
	})
}
