package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OkResponse struct {
	Success bool   `json:"success"`
	Info    string `json:"info"`
	Error   string `json:"error"`
}

func JsonOK(ctx *gin.Context, obj interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, obj)
}

// ErrorResponse response of error
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
