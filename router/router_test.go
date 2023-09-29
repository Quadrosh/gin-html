package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quadrosh/gin-html/config"
	"github.com/quadrosh/gin-html/controllers"
	"gotest.tools/v3/assert"
)

func TestPingRoute(t *testing.T) {
	var appConfig config.AppConfig
	appConfig.LoadConfig()

	db, err := config.ConnectDB(appConfig.LocalConfig)
	if err != nil {
		panic(err)
	}

	var ctl = &controllers.Controller{
		App: &appConfig,
		Db:  db,
	}

	var router = InitRoutes(ctl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}
