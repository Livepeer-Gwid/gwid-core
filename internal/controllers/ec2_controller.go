package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/services"
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
	params, exists := middleware.GetQueryParams(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get query params"})
		return
	}

	ec2InstanceTypes, statusCode, err := s.ec2Service.GetEC2InstanceTypes(params)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err,
		})

		return
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    ec2InstanceTypes,
	})
}
