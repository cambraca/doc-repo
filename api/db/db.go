package db

import (
	"api/models"
	"gorm.io/gorm"
	"log"
	"os"
)
import "gorm.io/driver/postgres"

var DB *gorm.DB

func Connect() error {
	var err error

	dsn := os.ExpandEnv("host=$DB_HOST port=$DB_PORT " +
		"user=$DB_USER password=$DB_PASSWORD " +
		"dbname=$DB_NAME sslmode=$DB_SSLMODE")

	DB, err = gorm.Open(postgres.Open(dsn))

	if err == nil {
		log.Print("Connected to Postgres!")
	}

	return err
}

func RunMigrations() error {
	log.Print("Running migrations...")
	err := DB.AutoMigrate(&models.Document{})
	if err == nil {
		log.Print("Migrations done!")
	}
	return err
}
