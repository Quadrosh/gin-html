package main

import (
	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/controllers"
	"github.com/quadrosh/gin-html/router"
)

func main() {

	var err error
	gin.SetMode(gin.ReleaseMode)

	var ctx = controllers.NewContext()
	var router = router.InitRoutes(ctx)

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
