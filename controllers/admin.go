package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/render"
)

// HomePage - admin home page
// @Summary admin home page
// @Description admin home page
// @ID AdminHomePage
// @Tags admin home page
// @Produce  html
// @Success 200  "Success"
// @Router /admin/ [GET]
func (ctl *Controller) AdminHomePage(ctx *gin.Context) {

	// identity, _ := r.Context().Value(middleware.ContextKeyIdentity).(*auth.Identity)
	// ctxDatabase, _ := r.Context().Value(middleware.ContextKeyDatabase).(*gorm.DB)

	// user := &repository.User{}

	render.AdminTemplate(ctl.App, ctl.Engine, ctx, "home.page.tmpl", ResponseMap{
		"title": "AdminHomePage()",
	})

	// ctx.HTML(http.StatusOK, "admin.home.page.tmpl", gin.H{
	// 	"title": "AdminHomePage()",
	// })
}
