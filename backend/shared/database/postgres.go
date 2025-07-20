// shared/database/postgres.go
package database

import (
	"certitrack/backend/feature/certifications"
	"fmt"
	"log"
	"os"
	"strconv"

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

func ConnectDB() {
	var err error

	// Obtener configuración de las variables de entorno o usar valores por defecto
	dbHost := getEnv("DB_HOST", defaultDBHost)
	dbPort := getEnvAsInt("DB_PORT", defaultDBPort)
	dbUser := getEnv("DB_USER", defaultDBUser)
	dbPassword := getEnv("DB_PASSWORD", defaultDBPassword)
	dbName := getEnv("DB_NAME", defaultDBName)

	// Construir cadena de conexión
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
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
