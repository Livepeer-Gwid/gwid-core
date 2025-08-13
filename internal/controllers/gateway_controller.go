package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/types"
)

type GatewayController struct {
	gatewayService *services.GatewayService
}

func NewGatewayController(gatewayService *services.GatewayService) *GatewayController {
	return &GatewayController{
		gatewayService: gatewayService,
	}
}

func (gc *GatewayController) CreateAWSGateway(c *gin.Context) {
	createAWSGatewayReq := c.MustGet("validatedInput").(types.CreateGatewayWithAWSReq)

	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

	gateway, statusCode, err := gc.gatewayService.CreateGatewayWithAWS(createAWSGatewayReq, reqUser.ID)

	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    gateway,
	})
}

func (gc *GatewayController) GetUserGateways(c *gin.Context) {
	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

	params, exists := middleware.GetQueryParams(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "error": "failed to get query params"})
		return
	}

	data, statusCode, err := gc.gatewayService.GetUserGateways(reqUser.ID, params)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	total, err := gc.gatewayService.GetUserGatewaysCount(reqUser.ID)
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
