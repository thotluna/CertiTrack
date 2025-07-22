package main

import (
	"certitrack/backend/feature/certifications"
	"certitrack/backend/feature/person"
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

	router.GET("/debug/routes", func(c *gin.Context) {
		routes := router.Routes()
		c.JSON(http.StatusOK, routes)
	})

	apiV1 := router.Group("/api/v1")
	{
		repo := person.NewGormRepository(database.DB)
		service := person.NewService(repo)
		handler := person.NewHandler(service)
		handler.RegisterRoutes(apiV1)

		certHandler := certifications.NewHandler(database.DB)
		certHandler.RegisterRoutes(apiV1)

		log.Println("Registered routes:")
		for _, r := range router.Routes() {
			log.Printf("%s %s\n", r.Method, r.Path)
		}
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
