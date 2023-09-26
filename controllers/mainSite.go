package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/render"
)

// HomePage - front home page
// @Summary home page
// @Description main site home page
// @ID HomePage
// @Tags home page
// @Produce  html
// @Success 200  "Success"
// @Router / [GET]
func (ctl *Controller) HomePage(ctx *gin.Context) {

	render.MainTemplate(ctl.App, ctl.Engine, ctx, "home.page.tmpl", ResponseMap{
		"title": "Main HomePage()",
	})

}
