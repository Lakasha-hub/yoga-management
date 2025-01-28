package middlewares

import (
	"net/http"
	"yoga-management/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateMiddleware(ctx *gin.Context) {
	// Verify if token exists
	tokenStr, err := ctx.Cookie("tkn")
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	// Verify if token is valid
	token, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		ctx.Abort()
		return
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Claims"})
		ctx.Abort()
		return
	}

	// Set Role in gin Ctx
	if role, exists := claims["role"].(string); exists {
		ctx.Set("role", role)
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
		ctx.Abort()
		return
	}

	ctx.Next()
}
