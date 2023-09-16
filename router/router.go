package router

import (
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

	router.LoadHTMLGlob("templates/***/**/*")

	// var files []string
	// filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
	// 	if strings.HasSuffix(path, ".tmpl") {
	// 		files = append(files, path)
	// 	}
	// 	return nil
	// })

	router.Static("/static", "./static/")

	router.GET("/ping", conx.Ping)
	router.GET("/", conx.HomePage)

	// adminRouter := router.Group("/admin")
	// adminRouter.GET("", conx.AdminIndex)

	return router
}
