package main

import (
	"log"
	"os"
	"yoga-management/internal/db"
	"yoga-management/internal/handlers"
	"yoga-management/internal/middlewares"
	"yoga-management/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Migrate Models of DB
	db.Database.AutoMigrate(&models.Class{})
	db.Database.AutoMigrate(&models.User{})

	// Router config
	router := gin.Default()

	// Public Paths
	router.GET("/register", handlers.CreateUser)
	router.GET("/login", handlers.Login)

	// Protected Paths
	protected := router.Group("/api", middlewares.AuthenticateMiddleware)
	{
		protected.GET("/classes", handlers.GetClasses)
		protected.GET("/classes/:id", handlers.GetClassByID)
		protected.POST("/classes", handlers.CreateClass)
		protected.PUT("/classes/:id", handlers.UpdateClass)
		protected.DELETE("/classes/:id", handlers.DeleteClass)

	}

	//Listen and Serve APP
	port := os.Getenv("APP_PORT")
	router.Run(port)

}
