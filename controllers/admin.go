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
	"github.com/quadrosh/gin-html/responses"
)

// AdminHomePageResponse response to AdminHomePage
type AdminHomePageResponse struct {
	responses.OkResponse

	PageMeta
	CSRFResponse
}

// AdminHomePage - admin home page
// @Summary admin home page
// @Description admin home page
// @ID AdminHomePage
// @Tags admin home page
// @Produce  html
// @Success 200  "Success"
// @Router /admin/ [GET]
func (ctl *Controller) AdminHomePage(ctx *gin.Context) {

	var identity = ctx.Keys[constants.ContextKeyIdentity].(*auth.Identity)
	var user = repository.User{}

	if err := user.GetByID(ctl.Db, identity.User.ID); err != nil {
		ctl.ErrorPage(ctx, http.StatusBadRequest, errors.New(resources.UserNotFound()))
	}

	render.AdminTemplate(ctl.App, ctl.Engine, ctx, "home.page.tmpl", &AdminHomePageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Info:    ctl.GetStringFromSession(ctx, constants.SessionKeyInfo, true),
		},

		PageMeta: PageMeta{
			Title: "AdminHomePage()" + user.FirstName + " " + user.LastName,
		},
	})

}
