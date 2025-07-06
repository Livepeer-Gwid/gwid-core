package router

import (
	"net/http"
	"regexp"
	"time"

	"github.com/dvwright/xss-mw"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gwid.io/gwid-core/internals/config"
	"gwid.io/gwid-core/internals/middleware"
)

func NewRouter(cfg *config.Config) *gin.Engine {
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
