package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/render"
	"github.com/quadrosh/gin-html/repository"
	resources "github.com/quadrosh/gin-html/resources/ru"
	"github.com/quadrosh/gin-html/responses"
	"gorm.io/gorm"
)

// PasswordResetPostResponse response to PasswordReset request
type PasswordResetPostResponse struct {
	responses.OkResponse
}

// PasswordResetPage - password reset page
// @Summary password reset page
// @Description page of user password reset
// @ID PasswordResetPage
// @Tags password reset
// @Produce  html
// @Success 200
// @Router /password-reset/:token [GET]
func (ctl *Controller) PasswordResetPage(ctx *gin.Context) {
	var token = ctx.Param("token")
	if token == "" {
		ctl.ErrorPage(ctx, http.StatusBadRequest, errors.New("token not found"))
		return
	}

	var (
		user = repository.User{}
		db   = ctl.Db
	)
	err := user.GetByPasswordResetToken(db, token)
	if err != nil {
		log.Println("PasswordResetPage() err: ", err.Error())
		if err == gorm.ErrRecordNotFound {
			ctl.ErrorPage(ctx, http.StatusBadRequest, errors.New(resources.LinkIsOld()))
			return
		}
		ctl.ErrorPage(ctx, http.StatusInternalServerError, err)
		return
	}

	render.MainTemplate(ctl.App, ctl.Engine, ctx, "password-reset.page.tmpl", ResponseMap{
		"title":      "AdminHomePage()",
		"token":      token,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})

}

// PasswordResetForm form of password reset
type PasswordResetForm struct {
	Password string `form:"password" binding:"required"`
}

// PasswordResetPOST - password reset request
// @Summary password reset
// @Description password reset
// @ID PasswordResetPOST
// @Tags password reset
// @Produce  json
// @Success 200
// @Router /password-reset/:token [POST]
func (ctl *Controller) PasswordResetPOST(ctx *gin.Context) {

	var token = ctx.Param("token")
	if token == "" {
		ctl.ErrorJSON(ctx, errors.New("token not found"))
		return
	}

	var form PasswordResetForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctl.ErrorJSON(ctx, err)
		return
	}

	var (
		user = repository.User{}
		db   = ctl.Db
	)
	err := user.GetByPasswordResetToken(db, token)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(resources.LinkIsOld()))
		return
	}

	if user.ID == 0 {
		ctl.ErrorJSON(ctx, errors.New(resources.LinkIsOld()))
		return
	}

	err = user.HashPassword(form.Password)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(resources.СantHashPassword()))
		return
	}

	if err := db.Model(&repository.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"password_hash":        user.PasswordHash,
			"password_reset_token": "",
		}).Error; err != nil {
		ctl.ErrorJSON(ctx, errors.New(resources.DatabaseSaveError()))
		return
	}

	responses.JsonOK(ctx, PasswordResetPostResponse{
		OkResponse: responses.OkResponse{Success: true},
	})
}
