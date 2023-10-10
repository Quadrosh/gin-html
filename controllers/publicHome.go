package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/internal/constants"
	"github.com/quadrosh/gin-html/render"
	"github.com/quadrosh/gin-html/repository"
	resources "github.com/quadrosh/gin-html/resources/ru"
	"github.com/quadrosh/gin-html/responses"
)

type PublicHomePageResponse struct {
	responses.OkResponse
	Page                  publicArticleEntry
	ArticleLayoutConstMap map[string]repository.ArticleLayout

	CSRFResponse
}

// PublicHomePage - front home page
// @Summary home page
// @Description main site home page
// @ID PublicHomePage
// @Tags home page
// @Produce  html
// @Success 200  "Success"
// @Router / [GET]
func (ctl *Controller) PublicHomePage(ctx *gin.Context) {

	var pageURL = "home"

	var (
		db   = ctl.Db
		page = repository.Article{}
	)

	err := page.ByURL(db, pageURL)
	if err != nil {
		ctl.ErrorPage(ctx, 404, errors.New(resources.PageNotFound()))
		return
	}

	var entry publicArticleEntry
	err = entry.convert(&page)
	if err != nil {
		ctl.ErrorPage(ctx, 500, err)
		return
	}

	render.PublicTemplate(ctl.App, ctl.Engine, ctx, "home.page.tmpl", &PublicHomePageResponse{
		OkResponse: responses.OkResponse{
			Success: true,
			Error:   ctl.GetSessionString(ctx, constants.SessionKeyError, true),
			Info:    ctl.GetSessionString(ctx, constants.SessionKeyInfo, true),
		},
		Page:                  entry,
		ArticleLayoutConstMap: repository.ArticleLayoutConstMap,
	})

}
