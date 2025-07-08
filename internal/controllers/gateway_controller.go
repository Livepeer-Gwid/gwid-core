package controllers

import (
	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/models"
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

func (gc *GatewayController) CreateGateway(c *gin.Context) {
	createGatewayReq := c.MustGet("validatedInput").(types.DeployGatewayPayloadReq)

	reqUser := c.MustGet("user").(*types.JwtCustomClaims)

	gateway := models.Gateway{
		Provider:           createGatewayReq.Provider,
		Region:             createGatewayReq.Region,
		GatewayName:        createGatewayReq.GatewayName,
		GatewayType:        createGatewayReq.GatewayType,
		RPCURL:             createGatewayReq.RPCURL,
		Password:           createGatewayReq.Password,
		TranscodingProfile: createGatewayReq.TranscodingProfile,
		Status:             models.Initializing,
		UserID:             reqUser.ID,
	}

	statusCode, err := gc.gatewayService.CreateGateway(&gateway)
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
