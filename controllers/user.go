package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/internal/auth"
	"github.com/quadrosh/gin-html/internal/constants"
	"github.com/quadrosh/gin-html/render"
	"github.com/quadrosh/gin-html/repository"
	resources "github.com/quadrosh/gin-html/resources/ru"
)

// UserHomePage - user home page
// @Summary user home page
// @Description user home page
// @ID UserHomePage
// @Tags user home page
// @Produce  html
// @Success 200  "Success"
// @Router /user/ [GET]
func (ctl *Controller) UserHomePage(ctx *gin.Context) {

	var identity = ctx.Keys[constants.ContextKeyIdentity].(*auth.Identity)
	var user = repository.User{}

	if err := user.GetByID(ctl.Db, identity.User.ID); err != nil {
		ctl.ErrorPage(ctx, http.StatusBadRequest, errors.New(resources.UserNotFound()))
	}

	render.UserTemplate(ctl.App, ctl.Engine, ctx, "home.page.tmpl", ResponseMap{
		"title": "UserHomePage " + user.FirstName + " " + user.LastName,
	})

}
