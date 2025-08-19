package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/services"
	"gwid.io/gwid-core/internal/types"
)

type RegionController struct {
	regionService *services.RegionService
}

func NewRegionController(regionService *services.RegionService) *RegionController {
	return &RegionController{
		regionService: regionService,
	}
}

func (s *RegionController) GetAWSRegions(c *gin.Context) {
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

	region, status, err := s.regionService.GetAWSRegions(reqUser.ID, credentialsID)
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	c.JSON(status, gin.H{
		"success": true,
		"data":    region,
	})
}
