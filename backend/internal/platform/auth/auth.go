package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func AuthenticateMiddleware(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Verify if token exists
		tokenStr, err := ctx.Cookie("tkn")
		if err != nil {
			log.Println(err)
			// ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}

		// Verify if token is valid
		token, err := ValidateJWT(tokenStr)
		if err != nil {
			log.Println(err)
			// ctx.Redirect(http.StatusSeeOther, "/login")
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

		// Verify if role is in roles
		if role, exists := claims["role"].(string); exists {
			roleExists := false
			for _, r := range roles {
				if r == role {
					roleExists = true
					ctx.Set("role", role)
					break
				}
			}
			if !roleExists {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Role not allowed"})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

var key = []byte(os.Getenv("JWT_KEY"))

// Generate JWT
func GenerateJWT(ctx *gin.Context, email, role string) error {
	log.Println("Generating JWT for:", email, role)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  email,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	result, err := token.SignedString(key)
	log.Println("Error signing token:", err)
	if err != nil {
		return err
	}
	log.Println("Token generated successfully:", result)
	ctx.SetCookie("tkn", result, 3600, "/", "", false, true)
	return nil
}

// Validate JWT
func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid Token")
	}

	return token, nil
}

func DeleteJWT(ctx *gin.Context) {
	ctx.SetCookie("tkn", "", -1, "/", "", false, true)
}

// Encrypts the sent password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Verify if the encrypted password is the same as the sent password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
