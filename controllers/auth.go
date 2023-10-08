package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/internal/auth"
	"github.com/quadrosh/gin-html/internal/constants"
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

// SignInRedirect тип редиректа
type SignInRedirect string

const (
	// SigninRedirectToUser redirect to user account page
	SigninRedirectToUser SignInRedirect = "/user"
	// SigninRedirectToAdmin redirect to admin account page
	SigninRedirectToAdmin SignInRedirect = "/admin"
)

// SigninResponse ответ на Signin
type SigninResponse struct {
	responses.OkResponse
	AccessToken string         `json:"access_token"`
	Redirect    SignInRedirect `json:"redirect"`
}

// SigninPageResponse response to SigninPage
type SigninPageResponse struct {
	responses.OkResponse
	responses.ConfirmResponse

	PageMeta
	CSRFResponse
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

	render.PublicTemplate(ctl.App, ctl.Engine, ctx, "password-reset.page.tmpl", ResponseMap{
		"title":      "AdminHomePage()", // TODO from pages
		"token":      token,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})

}

// PasswordResetForm form of password reset
type PasswordResetForm struct {
	Password string `form:"password" binding:"required"`
}

// SignInForm form of signin in
type SignInForm struct {
	Email    string `form:"email" binding:"required"`
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
		ctl.ErrorJSON(ctx, errors.New("token not found"), true)
		return
	}

	var form PasswordResetForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctl.ErrorJSON(ctx, err, true)
		return
	}

	var (
		user = repository.User{}
		db   = ctl.Db
	)
	err := user.GetByPasswordResetToken(db, token)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(resources.LinkIsOld()), true)
		return
	}

	if user.ID == 0 {
		ctl.ErrorJSON(ctx, errors.New(resources.LinkIsOld()), true)
		return
	}

	err = user.HashPassword(form.Password)
	if err != nil {
		ctl.ErrorJSON(ctx, errors.New(resources.СantHashPassword()), true)
		return
	}

	if err := db.Model(&repository.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"password_hash":        user.PasswordHash,
			"password_reset_token": "",
		}).Error; err != nil {
		ctl.ErrorJSON(ctx, errors.New(resources.DatabaseSaveError()), true)
		return
	}

	responses.JsonOK(ctx, PasswordResetPostResponse{
		OkResponse: responses.OkResponse{Success: true},
	})
}

// SigninPage - sign in page
// @Summary sign in page
// @Description sign in page
// @ID SigninPage
// @Tags signin
// @Produce  html
// @Success 200
// @Router /signin [GET]
func (ctl *Controller) SigninPage(ctx *gin.Context) {

	var sErr = ctl.GetSessionString(ctx, constants.SessionKeyError, true)

	if err := render.PublicTemplate(ctl.App, ctl.Engine, ctx, "signin.page.tmpl", &SigninPageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   sErr,
		},
		PageMeta: PageMeta{
			Title: "Sign in page()", // TODO from pages,
		},
	}); err != nil {
		log.Panic(err)
	}
}

// SigninPost - sign in form request
// @Summary  authorisation
// @Description  sign in form request
// @ID SigninPost
// @Tags  sign in
// @Produce  json
// @Success 200
// @Router /signin [POST]
func (ctl *Controller) SigninPost(ctx *gin.Context) {

	var form SignInForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctl.ErrorJSON(ctx, err, true)
		return
	}

	var (
		user = repository.User{}
		db   = ctl.Db
	)
	err := user.SignIn(db, form.Email, form.Password)
	if err != nil {
		ctl.ErrorJSON(ctx, err, true)
		return
	}

	if !user.RoleType.CanSignIn() {
		ctl.ErrorJSON(ctx, errors.New(resources.Forbidden()), true)
		return
	}

	var token string
	token, err = auth.CreateToken(user.ID, ctl.App.ApiTokenExpireSec, ctl.App.ApiSecret)
	if err != nil {
		ctl.ErrorJSON(ctx, err, true)
		return
	}

	// cохраняем токен в базу
	err = user.UpdateAuthKey(db, token)
	if err != nil {
		ctl.ErrorJSON(ctx, err, true)
		return
	}

	var redirect SignInRedirect
	if user.Can(db, repository.UserCanSettings{Rule: repository.UserRuleRoleAdmin}) {
		redirect = SigninRedirectToAdmin
	}
	if user.Can(db, repository.UserCanSettings{Rule: repository.UserRuleRoleUser}) {
		redirect = SigninRedirectToUser
	}

	if err := ctl.SetToSession(ctx, constants.SessionKeyInfo, resources.SignedInSuccessful()); err != nil {
		ctl.ErrorJSON(ctx, err, true)
		return
	}

	responses.JsonOK(ctx, SigninResponse{
		OkResponse:  responses.OkResponse{Success: true},
		AccessToken: token,
		Redirect:    redirect,
	})
}
