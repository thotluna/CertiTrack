// cmd/certitrack-backend/main.go
package main

import (
	"certitrack/backend/feature/certifications"
	"certitrack/backend/shared/database"
	"certitrack/backend/shared/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, assuming environment variables are set.")
	}

	database.ConnectDB()
	defer database.CloseDB()

	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(middleware.ErrorHandler())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from CertiTrack Backend with Gin and DB connection!",
		})
	})

	apiV1 := router.Group("/api/v1")
	{
		certHandler := certifications.NewHandler(database.DB)
		certHandler.RegisterRoutes(apiV1)
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}

	log.Printf("CertiTrack Backend running on port %s", apiPort)
	if err := router.Run(":" + apiPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
