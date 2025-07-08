// Package router sets up application router config
package router

import (
	"net/http"
	"regexp"
	"time"

	"github.com/dvwright/xss-mw"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internal/config"
	"gwid.io/gwid-core/internal/controllers"
	"gwid.io/gwid-core/internal/middleware"
	"gwid.io/gwid-core/internal/types"
)

func NewRouter(cfg *config.Config, authController *controllers.AuthController, userController *controllers.UserController, gatewayController *controllers.GatewayController) *gin.Engine {
	router := gin.Default()

	gin.SetMode(cfg.GinMode)

	setupRouteConfig(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"status":  "healthy",
			"message": "GWID core is running",
		})
	})

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/signup", middleware.ValidateRequestMiddleware[types.SignupReq](), authController.SignUp)
		auth.POST("/login", middleware.ValidateRequestMiddleware[types.LoginReq](), authController.Login)
	}

	user := router.Group("/api/v1/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile", userController.GetCurrentUserProfile)
	}

	gateway := router.Group("/api/v1/gateway")
	gateway.Use(middleware.AuthMiddleware())
	{
		gateway.POST("/", middleware.ValidateRequestMiddleware[types.DeployGatewayPayloadReq](), gatewayController.CreateGateway)
	}

	return router
}

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
