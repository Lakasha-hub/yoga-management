package handlers

import (
	"net/http"
	"os"
	"yoga-management/backend/internal/platform/auth"
	"yoga-management/backend/internal/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo user.UserRepository
}

func NewUserHandler(repo user.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	// Get json input values
	var json user.RegisterUserDTO
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(json.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}

	//Create new User
	newUser := user.User{
		NameUser: json.NameUser,
		Email:    json.Email,
		Password: hashedPassword,
	}

	result, err := h.repo.Register(&newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"payload": result})

}

func (h *UserHandler) LoginUser(ctx *gin.Context) {
	// Get json input values
	var json user.LoginUserDTO
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
		return
	}

	// Verify if user is ADMIN
	if json.Email == os.Getenv("ADMIN_EMAIL") && json.Password == os.Getenv("ADMIN_PASSWORD") {
		// Generate Admin cookie
		err := auth.GenerateJWT(ctx, json.Email, "admin")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"payload": "Login OK, Welcome Admin!"})
		return
	}

	// Verify if user exists
	user_exists, err := h.repo.Login(json.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Incorrect email or password"})
		return
	}

	// Verify password sent
	if err := auth.VerifyPassword(user_exists.Password, json.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect email or password"})
		return
	}

	// Set user cookie
	err = auth.GenerateJWT(ctx, json.Email, "user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"payload": "Login OK"})
}

func Logout(ctx *gin.Context) {
	auth.DeleteJWT(ctx)
	ctx.JSON(http.StatusOK, gin.H{"payload": "Log Out OK"})
}
