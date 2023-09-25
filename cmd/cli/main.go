package main

import (
	"log"
	"os"

	"github.com/quadrosh/gin-html/cli"
	"github.com/quadrosh/gin-html/config"
)

func main() {
	defer os.Exit(0)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var appConfig config.AppConfig
	appConfig.LoadConfig()

	db, err := config.ConnectDB(appConfig.LocalConfig)
	if err != nil {
		panic(err)
	}

	cmd := cli.CommandLine{}
	cmd.Run(db)
}
