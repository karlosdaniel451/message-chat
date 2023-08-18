package db

import (
	"fmt"
	"os"

	"github.com/karlosdaniel451/go-rest-api-template/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = os.Getenv("DB_HOST")
	user     = os.Getenv("DB_USER")
	port     = os.Getenv("DB_PORT")
	name     = os.Getenv("DB_NAME")
	password = os.Getenv("DB_PASSWORD")
)

var DB *gorm.DB

func Connect() error {
	var err error

	connectionConfig := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, name, port,
	)

	DB, err = gorm.Open(postgres.Open(connectionConfig), &gorm.Config{})

	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&model.Task{}) // Setup other models here
	if err != nil {
		return err
	}

	return nil
}
