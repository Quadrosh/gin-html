package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomePage - front home page
// @Summary home page
// @Description main site home page
// @ID HomePage
// @Tags home page
// @Produce  html
// @Success 200  "Success"
// @Router / [GET]
func (c *Context) HomePage(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "main.home.page.tmpl", gin.H{
		"title": "Main HomePage()",
	})
}
