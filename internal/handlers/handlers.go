package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"yoga-management/internal/db"
	"yoga-management/internal/models"
	"yoga-management/internal/utils"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
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

func CreateUser(ctx *gin.Context) {
	// Get json input values
	var json models.CreateUserDTO
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
		return
	}

	// Verify if user already exists
	var user_exists models.User
	if err := db.Database.Where("email = ?", json.Email).First(&user_exists).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(json.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}

	//Create new User
	newUser := models.User{
		NameUser: json.NameUser,
		Email:    json.Email,
		Password: hashedPassword,
	}

	// Save User in DB
	if err := db.Database.Save(&newUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user in database"})
		return
	}

	// Parse user info to response
	responseUser := models.ResponseUser{
		ID:        newUser.ID,
		NameUser:  newUser.NameUser,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{"payload": responseUser})

}

func Login(ctx *gin.Context) {
	// Verify if cookie already exists
	cookie, err := ctx.Cookie("tkn")
	if cookie != "" && err == nil {
		// Validate cookie
		_, err = utils.ValidateJWT(cookie)
		if err != nil {
			// Redirect to home
			ctx.JSON(http.StatusOK, gin.H{"payload": "Already logged in"})
			return
		}
	}

	// Get json input values
	var json models.LoginUserDTO
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
		return
	}

	// Token to set cookie
	var tokenStr string

	// Verify if user is ADMIN
	if json.Email == os.Getenv("ADMIN_EMAIL") && json.Password == os.Getenv("ADMIN_PASSWORD") {
		// Generate Admin cookie
		tokenStr, err := utils.GenerateJWT(json.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
			return
		}
		ctx.SetCookie("tkn", tokenStr, 3600, "/", "", false, true)
		ctx.JSON(http.StatusOK, gin.H{"payload": "Login OK, Welcome Admin!"})
		return
	}

	// Verify if user already exists
	var user_exists models.User
	if err := db.Database.Where("email = ?", json.Email).First(&user_exists).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Incorrect email or password"})
		return
	}

	// Verify password sent
	if err := utils.VerifyPassword(user_exists.Password, json.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or password"})
		return
	}

	// Set user cookie
	tokenStr, err = utils.GenerateJWT(json.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
		return
	}
	ctx.SetCookie("tkn", tokenStr, 3600, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"payload": "Login OK"})
}
