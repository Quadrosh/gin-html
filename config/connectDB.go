package config

import (
	"fmt"
	"log"

	"github.com/quadrosh/gin-html/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(lConfig LocalConfig) (*gorm.DB, error) {

	if !lConfig.InProduction {
		log.Println("Connecting to db...")
	}

	connString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		lConfig.DbHost,
		lConfig.DbPort,
		lConfig.DbName,
		lConfig.DbUser,
		lConfig.DbPass,
		lConfig.DbSSL)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if !lConfig.InProduction {
		log.Println("db connected")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return tx.AutoMigrate(repository.AllModels...)
	})
	if err != nil {
		panic(err)
	}

	return db, nil
}
