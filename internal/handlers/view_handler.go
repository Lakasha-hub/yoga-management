package handlers

import (
	"html/template"
	"net/http"
	"yoga-management/internal/db"
	"yoga-management/internal/models"

	"github.com/gin-gonic/gin"
)

const templateDir = "./templates/"

func Login(ctx *gin.Context) {
	data := struct {
		Title string
		Msg   string
	}{
		Title: "Yoga Center - Log in",
		Msg:   "Welcome to Yoga Center - Please Log in",
	}
	renderTemplate(ctx, "login.html", data)
}

func Register(ctx *gin.Context) {
	data := struct {
		Title string
		Msg   string
	}{
		Title: "Yoga Center - Sign up",
		Msg:   "Welcome to Yoga Center - Please Sign up",
	}
	renderTemplate(ctx, "register.html", data)
}

func Home(ctx *gin.Context) {

	// Get Classes
	var classes []models.Class
	db.Database.Find(&classes)

	data := struct {
		Classes []models.Class
	}{
		Classes: classes,
	}

	renderTemplate(ctx, "index.html", data)
}

func renderTemplate(ctx *gin.Context, page string, data any) {
	t := template.Must(template.ParseFiles("./templates/base.html", templateDir+page))

	err := t.ExecuteTemplate(ctx.Writer, "base", data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error when try to execute template"})
		return
	}
}
