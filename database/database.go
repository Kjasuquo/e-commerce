package database

import (
	"fmt"
	"github.com/decadevs/shoparena/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var Db *gorm.DB

func SetupDB() error {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, user, password, dbName, port)
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if db == nil {
		return fmt.Errorf("database was not initialized")
	}

	err = db.AutoMigrate(&models.Product{}, &models.Seller{}, &models.Buyer{}, &models.Order{}, models.Cart{})
	if err != nil {
		return fmt.Errorf("migration error: %v", err)
	}

	Db = db

	return nil
}
