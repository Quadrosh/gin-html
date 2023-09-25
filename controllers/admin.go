package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomePage - admin home page
// @Summary admin home page
// @Description admin home page
// @ID AdminHomePage
// @Tags admin home page
// @Produce  html
// @Success 200  "Success"
// @Router /admin/ [GET]
func (c *Context) AdminHomePage(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "admin.home.page.tmpl", gin.H{
		"title": "AdminHomePage()",
	})
}
