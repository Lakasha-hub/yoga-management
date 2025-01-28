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
	public := router.Group("/api")
	public.POST("/register", handlers.CreateUser)
	public.POST("/login", handlers.LoginUser)
	public.POST("/logout", handlers.Logout)
	public.GET("/classes", handlers.GetClasses)
	public.GET("/classes/:id", handlers.GetClassByID)

	// Views
	router.GET("/login", handlers.Login)
	router.GET("/unauthorized", handlers.Unauthorized)

	// Protected Paths
	protected := router.Group("/api", middlewares.AuthenticateMiddleware)
	protected.POST("/classes", handlers.CreateClass)
	protected.PUT("/classes/:id", handlers.UpdateClass)
	protected.DELETE("/classes/:id", handlers.DeleteClass)

	//Listen and Serve APP
	port := os.Getenv("APP_PORT")
	router.Run(port)

}
