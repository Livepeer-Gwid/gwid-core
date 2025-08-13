package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/models"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/types"
)

type AWSCredentialsController struct {
	awsCredentialsService *services.AWSCredentialsService
}

func NewAWSCredentialsController(awsCredentialsService *services.AWSCredentialsService) *AWSCredentialsController {
	return &AWSCredentialsController{
		awsCredentialsService: awsCredentialsService,
	}
}

func (ac *AWSCredentialsController) CreateAWSCredentials(c *gin.Context) {
	awsCredentialsReq := c.MustGet("validatedInput").(types.AWSCredentialsReq)

	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

	awsCredentials := models.AWSCredentials{
		AccessKeyID:     awsCredentialsReq.AccessKeyID,
		SecretAccessKey: awsCredentialsReq.SecretAccessKey,
		UserID:          reqUser.ID,
		RoleName:        "",
		RoleARN:         "",
		ProfieName:      "",
		ProfileARN:      "",
	}

	statusCode, err := ac.awsCredentialsService.CreateAWSCredentials(&awsCredentials)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	awsCredentials.User = nil

	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    awsCredentials,
	})
}

func (ac *AWSCredentialsController) GetUserAWSCredentials(c *gin.Context) {
	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

	params, exists := middleware.GetQueryParams(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get query params"})
		return
	}

	data, statusCode, err := ac.awsCredentialsService.GetUserAWSCredentials(reqUser.ID, params)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	total, err := ac.awsCredentialsService.GetUserCredentialsCount(reqUser.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	metadata := &types.Metadata{
		Total:  total,
		Count:  len(*data),
		Page:   params.Page,
		Limit:  params.Limit,
		Order:  params.Order,
		Search: params.Search,
	}

	c.JSON(statusCode, gin.H{
		"success":  true,
		"data":     data,
		"metadata": metadata,
	})
}
