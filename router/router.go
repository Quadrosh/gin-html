package router

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/controllers"
	"github.com/quadrosh/gin-html/internal/auth"
	"github.com/quadrosh/gin-html/internal/constants"
	"github.com/quadrosh/gin-html/repository"
	resources "github.com/quadrosh/gin-html/resources/ru"
	"github.com/quadrosh/gin-html/responses"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"
)

// InitRoutes инициализация доступных URL
func InitRoutes(ctl *controllers.Controller) *gin.Engine {

	var router = gin.Default() // with logger
	// router := gin.New()
	// router.Use(gin.Recovery())

	router.Use(secure.New(secure.Config{
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy: "strict-origin-when-cross-origin",
	}))

	var store = cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: ctl.App.ApiSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	router.Use(CORSMiddleware())
	router.Use(gin.Recovery())
	// router.Use(middleware.ErrorHandler())

	// router.Use(favicon.New("./favicon.ico"))

	router.Static("/static", "./static/")

	router.GET("/ping", ctl.Ping)
	router.GET("/", ctl.HomePage)
	router.GET("/home", ctl.HomePage)
	router.GET("/password-reset/:token", ctl.PasswordResetPage)
	router.POST("/password-reset-post/:token", ctl.PasswordResetPOST)

	router.GET(constants.RedirectErrURL, func(c *gin.Context) {
		var session = sessions.Default(c)
		sVal := session.Get(constants.ContextKeyRedirectError)
		if sVal == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": resources.SystemError()})
			return
		}
		var errResp, ok = sVal.(responses.ErrorResponse)
		if !ok || &errResp == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": resources.SystemError()})
			return
		}
		ctl.ErrorPage(c, errResp.Code, errors.New(errResp.Message))
		return
	})
	router.GET("/signin", ctl.SigninPage)
	router.POST("/signin", ctl.SigninPost)

	adminRoute := router.Group("/admin")
	adminRoute.Use(AuthMiddleware(ctl.Db, ctl.App.ApiSecret, repository.UserRoleTypeAdmin, router))
	adminRoute.GET("", ctl.AdminHomePage)
	adminRoute.GET("/pages", ctl.AdminPageIndexPage)

	userRoute := router.Group("/user")
	userRoute.Use(AuthMiddleware(ctl.Db, ctl.App.ApiSecret, repository.UserRoleTypeUser, router))
	userRoute.GET("", ctl.UserHomePage)

	ctl.Engine = router
	return router
}

// CORSMiddleware cors headers
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

// AuthMiddleware - authentification
func AuthMiddleware(db *gorm.DB, apiSecret string, role repository.UserRoleType, r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.CheckUser(
			c.Request,
			role,
			db,
			apiSecret,
		)
		if err != nil {
			c.Error(err)
			if c.Request.Method == "GET" {
				var session = sessions.Default(c)
				session.Set(constants.ContextKeyRedirectError, responses.ErrorResponse{
					Code:    http.StatusForbidden,
					Message: resources.Forbidden(),
				})
				c.Request.URL.Path = constants.RedirectErrURL
				r.HandleContext(c)
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": resources.Forbidden()})
			return
		}
		c.Set(constants.ContextKeyIdentity, &auth.Identity{User: user, IsAuthorized: true})
		c.Next()
	}
}
