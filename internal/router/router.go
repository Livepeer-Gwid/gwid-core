package router

import (
	"net/http"
	"regexp"
	"time"

	"github.com/dvwright/xss-mw"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/di"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/types"
)

func setupRouteConfig(router *gin.Engine) {
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			reg := regexp.MustCompile(`^(https://([a-zA-Z0-9-]+\.)*gwid\.io|http://localhost(:[0-9]+)?)$`)
			return reg.MatchString(origin)
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(middleware.RateLimitMiddleware())

	var xssMdlwr xss.XssMw
	router.Use(xssMdlwr.RemoveXss())
}

func SetupRoutes(container *di.Container) *gin.Engine {
	router := gin.Default()

	setupRouteConfig(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "pong",
		})
	})

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/signup", middleware.ValidateRequestMiddleware[types.SignupReq](), container.AuthController.SignUp)
		auth.POST("/login", middleware.ValidateRequestMiddleware[types.LoginReq](), container.AuthController.Login)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "route not found",
		})
	})

	return router
}
