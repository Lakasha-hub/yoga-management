package handlers

import (
	"net/http"
	"os"
	"yoga-management/internal/db"
	"yoga-management/internal/models"
	"yoga-management/internal/utils"

	"github.com/gin-gonic/gin"
)

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

func LoginUser(ctx *gin.Context) {
	// Verify if cookie already exists
	cookie, err := ctx.Cookie("tkn")
	if cookie != "" && err == nil {
		// Validate cookie
		_, err = utils.ValidateJWT(cookie)
		if err != nil {
			// Redirect to home
			//...
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

func Logout(ctx *gin.Context) {
	ctx.SetCookie("tkn", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"payload": "Log Out OK"})
}
