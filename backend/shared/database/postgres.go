package database

import (
	"certitrack/backend/feature/certifications"
	"certitrack/backend/feature/person"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultDBHost     = "localhost"
	defaultDBPort     = 5432
	defaultDBUser     = "user"
	defaultDBName     = "certitrack_db"
	defaultDBPassword = ""
)

var DB *gorm.DB

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func ConnectForTesting(connStr string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			tx := db.Exec("SELECT 1")
			if tx.Error == nil {
				break
			}
		}

		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to connect to test database after %d attempts: %v", maxRetries, err)
		}

		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %v", err)
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		return nil, fmt.Errorf("failed to create pgcrypto extension: %v", err)
	}

	if err := db.AutoMigrate(&person.Person{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %v", err)
	}

	return db, nil
}

func ConnectDB() {
	var err error

	dbHost := getEnv("DB_HOST", defaultDBHost)
	dbPort := getEnvAsInt("DB_PORT", defaultDBPort)
	dbUser := getEnv("DB_USER", defaultDBUser)
	dbPassword := getEnv("DB_PASSWORD", defaultDBPassword)
	dbName := getEnv("DB_NAME", defaultDBName)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	err = DB.AutoMigrate(
		&person.Person{},
		&certifications.Certification{},
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
