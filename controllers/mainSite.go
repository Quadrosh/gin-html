package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomePage - front home page
// @Summary ping-pong test
// @Description test of working server
// @ID Ping
// @Tags ping-pong
// @Produce  json
// @Success 200  "Success"
// @Router /ping [GET]
func (c *Context) HomePage(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "home.page.tmpl", gin.H{
		"title": "Main website",
	})
}
