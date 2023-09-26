package router

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/controllers"
	csrf "github.com/utrack/gin-csrf"
)

// InitRoutes инициализация доступных URL
func InitRoutes(ctl *controllers.Controller) *gin.Engine {

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

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	router.Use(CORSMiddleware())
	// router.Use(middleware.ErrorHandler())

	// router.Use(favicon.New("./favicon.ico"))

	router.Static("/static", "./static/")

	router.GET("/ping", ctl.Ping)
	router.GET("/", ctl.HomePage)
	router.GET("/password-reset/:token", ctl.PasswordResetPage)
	router.POST("/password-reset-post/:token", ctl.PasswordResetPOST)

	adminRouter := router.Group("/admin")
	adminRouter.GET("", ctl.AdminHomePage)

	ctl.Engine = router
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
