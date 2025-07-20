// shared/database/postgres.go
package database

import (
	"certitrack/backend/feature/certifications"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// Import other feature models that need migration here
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// Build the connection string using environment variables
	// Removed TimeZone=America/Caracas as it's not recognized by Alpine Linux
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Configure GORM log level
	})
	if err != nil {
		// Changed log.Fatalf to log.Panicf to ensure stack trace is printed
		log.Panicf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	// Execute auto migrations
	err = DB.AutoMigrate(
		&certifications.Certification{},
		// Add other feature models that need migration here
	)
	if err != nil {
		log.Panicf("Failed to auto migrate database: %v", err)
	}

	log.Println("Database migrations completed successfully!")
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
	log.Println("Database connection closed.")
}
