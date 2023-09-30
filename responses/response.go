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

// ConfirmResponse response of confirm modal message
type ConfirmResponse struct {
	Message      string `json:"message"`
	YesBtnName   string `json:"yes_btn_name"`
	YesBtnAction string `json:"yes_btn_action"`
	NoBtnName    string `json:"no_btn_name"`
}
