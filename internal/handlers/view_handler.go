package handlers

import (
	"html/template"
	"net/http"

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

func renderTemplate(ctx *gin.Context, page string, data any) {
	t := template.Must(template.ParseFiles("./templates/base.html", templateDir+page))

	err := t.ExecuteTemplate(ctx.Writer, "base", data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error when try to execute template"})
		return
	}
}
