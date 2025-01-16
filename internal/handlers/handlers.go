package handlers

import (
	"fmt"
	"net/http"
	"time"
	"yoga-management/internal/db"
	"yoga-management/internal/models"

	"github.com/gin-gonic/gin"
)

func GetClasses(ctx *gin.Context) {
	var classes []models.Class
	db.Database.Find(&classes)
	ctx.JSON(http.StatusOK, gin.H{"payload": classes})
}

func GetClassByID(ctx *gin.Context) {
	// Get id of class
	id := ctx.Param("id")

	// Find first Class with ID sent
	class := models.Class{}
	db.Database.First(&class, id)
	if class.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("There is no Class with id %s", id)})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"payload": class})
}

func CreateClass(ctx *gin.Context) {
	// Get json input values
	var json models.CreateClassDTO
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse DateTime
	dateClass, err := time.Parse(time.DateTime, json.DateClass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date/time format. Use YYYYY-MM-DD HH:MM:SS"})
		return
	}

	// Verify if class already exists
	var class_exists models.Class
	if err := db.Database.Where("name_class = ? AND professor = ? AND date_class = ?", json.NameClass, json.Professor, dateClass).First(&class_exists).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Class already exists"})
		return
	}

	// Parse CreateClassDTO to DB model Class
	newClass := models.Class{
		NameClass: json.NameClass,
		Professor: json.Professor,
		DateClass: dateClass,
		Capacity:  json.Capacity,
	}

	// Save new Class in DB
	if err := db.Database.Save(&newClass).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating class in database"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"payload": newClass})
}

func UpdateClass(ctx *gin.Context) {

	// Get id of class
	id := ctx.Param("id")

	// Verify Class exists
	var class models.Class
	if err := db.Database.First(&class, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Clase no encontrada"})
		return
	}

	// Get json input values
	var json models.UpdateClassDTO
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only update non-zero fields by default
	updatedValues := models.Class{
		NameClass: json.NameClass,
		Professor: json.Professor,
		Capacity:  json.Capacity,
	}

	// Parse DateTime if sent
	if dateClass, err := time.Parse(time.DateTime, json.DateClass); json.DateClass != "" {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date/time format. Use YYYYY-MM-DD HH:MM:SS"})
			return
		}
		updatedValues.DateClass = dateClass
	}

	db.Database.Model(&class).Updates(updatedValues)
	ctx.JSON(http.StatusAccepted, gin.H{"payload": class})
}

func DeleteClass(ctx *gin.Context) {
	// Get id of class
	id := ctx.Param("id")

	// Find first Class with ID sent
	class := models.Class{}
	db.Database.First(&class, id)
	if class.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("There is no Class with id %s", id)})
		return
	}

	// Delete Class by ID
	db.Database.Delete(&class)
	ctx.JSON(http.StatusOK, gin.H{"payload": "The Class has been deleted succesfully"})
}
