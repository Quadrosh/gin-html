package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/render"
	"github.com/quadrosh/gin-html/responses"
)

// ErrorPage - page of error
// @Summary error page
// @Description error page
// @ID ErrorPage
// @Tags error page
// @Produce  html
// @Success 200  "Success"
// @Router /error [GET]
func (c *Controller) ErrorPage(ctx *gin.Context, code int, e error) {

	rErr := render.MainTemplate(c.App, c.Engine, ctx, "error.page.tmpl", ResponseMap{
		"title":       "Error",
		"description": "During execution occurs error",
		"code":        code,
		"error":       e.Error(),
	})
	if rErr != nil {
		panic(rErr)
	}
}

// ErrorJSON - error response
// @Summary error json
// @Description error json response
// @ID ErrorJSON
// @Tags error json
// @Produce  json
func (c *Controller) ErrorJSON(ctx *gin.Context, err error) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, responses.OkResponse{
		Success: false,
		Error:   err.Error(),
	})

}
