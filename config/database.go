package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	logMode := logger.Silent
	if os.Getenv("ENV") == "development" {
		logMode = logger.Info
	}

	for i := 0; i < 5; i++ {
		fmt.Printf("Attempting to connect to database (attempt %d)...\n", i+1)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logMode),
		})

		if err == nil {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(30 * time.Minute)

			fmt.Println("Database connected successfully!")
			DB = db
			return
		}

		fmt.Println("Database connection failed:", err)
		time.Sleep(3 * time.Second)
	}

	log.Fatal("Database connection failed after multiple attempts")
}
