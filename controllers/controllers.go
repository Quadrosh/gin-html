package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/config"
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

// Ping - ping-pong test
// @Summary ping-pong test
// @Description test of working server
// @ID Ping
// @Tags ping-pong
// @Produce  json
// @Success 200  "Success"
// @Router /ping [GET]
func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
