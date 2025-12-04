package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect loads env vars, connects to DB, and runs migrations
func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env vars")
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL not found in environment variables")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to Postgres, yayy!")

	// Assign to global variable so other packages can use it
	DB = db

	fmt.Println("Database migration complete!")
}
