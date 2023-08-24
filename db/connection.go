package db

import (
	"fmt"
	"os"

	"github.com/karlosdaniel451/message-chat/domain/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	host     string
	user     string
	port     string
	name     string
	password string
)

var DB *gorm.DB

func Connect() error {
	host = os.Getenv("DB_HOST")
	user = os.Getenv("DB_USER")
	port = os.Getenv("DB_PORT")
	name = os.Getenv("DB_NAME")
	password = os.Getenv("DB_PASSWORD")

	var err error

	connectionConfig := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, name, port,
	)

	DB, err = gorm.Open(
		postgres.Open(connectionConfig), 
		&gorm.Config{Logger: logger.Discard},
	)

	if err != nil {
		return err
	}

	err = DB.AutoMigrate(
		&model.Group{},
		&model.GroupMessage{},
		&model.PrivateMessage{},
		&model.User{},
		// Setup other models here
	)

	if err != nil {
		return err
	}

	return nil
}
