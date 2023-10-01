package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/config"
	"github.com/quadrosh/gin-html/internal/constants"
	"github.com/quadrosh/gin-html/repository"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"
)

// Controller - контекст контроллера
type Controller struct {
	App    *config.AppConfig
	Db     *gorm.DB
	Engine *gin.Engine
}

// CSRFResponse struct containing CSRFToken data
type CSRFResponse struct {
	CSRFToken  string
	CurrentURL string
}

// CSRF implementation iCSRF
func (res *CSRFResponse) CSRF(ctx *gin.Context) {
	res.CSRFToken = csrf.GetToken(ctx)
	res.CurrentURL = ctx.Request.URL.Path
}

// ResponseMap response data map
type ResponseMap gin.H

// CSRF implementation iCSRF
func (m ResponseMap) CSRF(ctx *gin.Context) {
	m["CSRFToken"] = csrf.GetToken(ctx)
	m["CurrentURL"] = ctx.Request.URL.Path
}

// GetPagination - получение пагинации из GET
func (ctl *Controller) GetPagination(r *http.Request) repository.Pagination {
	var (
		result     repository.Pagination
		curPageStr = r.URL.Query().Get("page")
		count      = r.URL.Query().Get("size")
	)
	if c, err := strconv.Atoi(count); err == nil {
		result.PageSize = c
	} else {
		result.PageSize = constants.DefaultEntriesCount
	}
	if p, err := strconv.Atoi(curPageStr); err == nil {
		result.CurrentPage = p
	} else {
		result.CurrentPage = 1
	}

	return result
}

// GetStringFromSession gets string value by key from session and clear it
func (ctl *Controller) GetStringFromSession(ctx *gin.Context, key string, clearValue bool) string {
	var session = sessions.Default(ctx)
	value, _ := session.Get(key).(string)
	if value != "" && clearValue {
		session.Set(key, nil)
		session.Save()
	}
	return value
}

// SetToSession set value to session
func (ctl *Controller) SetToSession(ctx *gin.Context, key string, value interface{}) {
	var session = sessions.Default(ctx)
	session.Set(key, value)
	session.Save()
}

// Ping - ping-pong test
// @Summary ping-pong test
// @Description test of working server
// @ID Ping
// @Tags ping-pong
// @Produce  json
// @Success 200  "Success"
// @Router /ping [GET]
func (ctl *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
