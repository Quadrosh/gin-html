package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quadrosh/gin-html/config"
	"github.com/quadrosh/gin-html/controllers"
	"github.com/quadrosh/gin-html/router"
	"gorm.io/gorm"
)

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var appConfig config.AppConfig
	appConfig.LoadConfig()

	db, err := config.ConnectDB(appConfig.LocalConfig)
	if err != nil {
		panic(err)
	}

	RunServer(db, &appConfig)
}

func RunServer(
	db *gorm.DB,
	app *config.AppConfig,
) {

	var ctl = &controllers.Controller{
		App: app,
		Db:  db,
	}

	gin.SetMode(gin.ReleaseMode)

	var router = router.InitRoutes(ctl)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", app.AppPort), // host:port
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Minute,
		Handler:           router,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen error: %s\n", err)
		}
	}()

	log.Printf("app started on port %s \n", app.AppPort)
	log.Printf("dev server runs on http://localhost:%s \n", app.AppPort)

	<-shutdownChan
	log.Println("shutting down")

	_ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// other handling
		cancel()
	}()
	if err := server.Shutdown(_ctx); err != nil {
		log.Panicf("server shutdown failed: %+v", err)
	}

	log.Println("shutted down gracefully")
}
