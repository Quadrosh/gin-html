package router

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/controllers"
)

// InitRoutes инициализация доступных URL
func InitRoutes(conx *controllers.Context) *gin.Engine {

	ginMode := "release"
	gin.SetMode(ginMode)

	var router = gin.Default() // with logger
	// router := gin.New()
	// router.Use(gin.Recovery())

	router.Use(secure.New(secure.Config{
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy: "strict-origin-when-cross-origin",
	}))
	router.Use(CORSMiddleware())

	// router.Use(favicon.New("./favicon.ico"))

	router.LoadHTMLGlob("templates/**/*")

	router.Static("/static", "./static/")

	router.GET("/ping", conx.Ping)
	router.GET("/", conx.HomePage)

	adminRouter := router.Group("/admin")
	adminRouter.GET("", conx.AdminHomePage)

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
