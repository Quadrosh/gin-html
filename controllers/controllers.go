package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping - ping-pong test
// @Summary ping-pong test
// @Description test of working server
// @ID Ping
// @Tags ping-pong
// @Produce  json
// @Success 200  "Success"
// @Router /ping [GET]
func (c *Context) Ping(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
