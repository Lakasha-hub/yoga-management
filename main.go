package main

import (
	"log"
	"net/http"
	"os"
	"yoga-management/internal/db"
	"yoga-management/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Database.AutoMigrate(&models.Class{})

	router := gin.Default()

	router.GET("/classes", func(ctx *gin.Context) {
		var classes []models.Class
		db.Database.Find(&classes)
		ctx.JSON(http.StatusOK, classes)
	})

	port := os.Getenv("APP_PORT")
	router.Run(port)

}
