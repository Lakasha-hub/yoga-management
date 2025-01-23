package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var key = []byte(os.Getenv("JWT_KEY"))

// Generate JWT
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  email,
		"role": getRole(email),
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	result, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return result, nil
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

func getRole(email string) string {
	if email == os.Getenv("ADMIN_EMAIL") {
		return "admin"
	}
	return "user"
}
