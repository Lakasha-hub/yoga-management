package main

import (
	"log"
	"os"
	"yoga-management/internal/db"
	"yoga-management/internal/handlers"
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

	// Router config
	router := gin.Default()

	// Paths
	router.GET("/classes", handlers.GetClasses)
	router.GET("/classes/:id", handlers.GetClassByID)
	router.POST("/classes", handlers.CreateClass)
	router.PUT("/classes/:id", handlers.UpdateClass)
	router.DELETE("/classes/:id", handlers.DeleteClass)

	//Listen and Serve APP
	port := os.Getenv("APP_PORT")
	router.Run(port)

}
