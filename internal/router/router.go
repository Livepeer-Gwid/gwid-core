// Package router sets up application router config
package router

import (
	"log"
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

func NewRouter(
	cfg *config.Config,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	gatewayController *controllers.GatewayController,
	regionController *controllers.RegionController,
	awsCredentialsController *controllers.AWSCredentialsController,
	ec2Controller *controllers.EC2Controller,
) *gin.Engine {
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
		user.GET("/gateway", middleware.QueryMiddleware(), gatewayController.GetUserGateways)
	}

	gateway := router.Group("/api/v1/gateway")
	gateway.Use(middleware.AuthMiddleware())
	{
		gateway.POST("/aws", middleware.ValidateRequestMiddleware[types.CreateGatewayWithAWSReq](), gatewayController.CreateAWSGateway)
	}

	region := router.Group("/api/v1/region")
	region.Use(middleware.AuthMiddleware())
	{
		region.GET("/aws", middleware.QueryMiddleware(), regionController.GetAWSRegions)
	}

	awsCredentials := router.Group("/api/v1/aws-credentials")
	awsCredentials.Use(middleware.AuthMiddleware())
	{
		awsCredentials.POST("/", middleware.ValidateRequestMiddleware[types.AWSCredentialsReq](), awsCredentialsController.CreateAWSCredentials)
		awsCredentials.GET("/", middleware.QueryMiddleware(), awsCredentialsController.GetUserAWSCredentials)
	}

	ec2 := router.Group("/api/v1/ec2")
	ec2.Use(middleware.AuthMiddleware())
	{
		ec2.GET("/", middleware.QueryMiddleware(), ec2Controller.GetEC2InstanceTypes)
	}

	return router
}

func setupRouteConfig(router *gin.Engine) {
	err := router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	if err != nil {
		log.Fatalln("router not initiated")
	}

	originRegex := regexp.MustCompile(`^https?:\/\/(localhost(:\d+)?|([a-zA-Z0-9-]+\.)?gwid\.io)$`)

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return originRegex.MatchString(origin)
		},
		// AllowOrigins: []string{"https://gwid.io", "https://app.gwid.io", "http://localhost:3000"},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(middleware.RateLimitMiddleware())

	var xssMdlwr xss.XssMw
	router.Use(xssMdlwr.RemoveXss())
}
