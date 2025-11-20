package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	_ = godotenv.Load()
	// If a full DATABASE_URL is provided (e.g. by Render), prefer that.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to the Database", err)
		}

		DB = db
		return
	}
	host := os.Getenv("HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=%s", host, user, port, password, name, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database", err)
	}

	DB = db
}
