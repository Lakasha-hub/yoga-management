package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

const templateDir = "./templates/"

func Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Login Page"})
	// renderTemplate(ctx, "login.html", nil)
}

func Register(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Register Page"})
	// renderTemplate(ctx, "register.html", nil)
}

func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Home Page"})
	// renderTemplate(ctx, "index.html", nil)
}

func renderTemplate(ctx *gin.Context, page string, data any) {
	t := template.Must(template.ParseFiles("../../../frontend/templates/base.html", templateDir+page))

	err := t.ExecuteTemplate(ctx.Writer, "base", data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error when try to execute template"})
		return
	}
}
