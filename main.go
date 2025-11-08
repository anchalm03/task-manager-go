package main

import (
	"fmt"
	"log"
	"os"

	"task_manager/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env vars")
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL not found in environment variables")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET not set in .env")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to Postgres, yayy!")

	// Auto-migrate all models
	err = db.AutoMigrate(&models.User{}, &models.Project{}, &models.Task{}, &models.Comment{})
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}
	fmt.Println("Database migration complete!")

	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
