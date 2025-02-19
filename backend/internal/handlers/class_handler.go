package handlers

import (
	"net/http"
	"strconv"
	"time"
	"yoga-management/backend/internal/class"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClassHandler struct {
	repo class.ClassRepository
}

func NewClassHandler(repo class.ClassRepository) *ClassHandler {
	return &ClassHandler{repo: repo}
}

func (h *ClassHandler) GetClasses(ctx *gin.Context) {
	result, err := h.repo.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting classes"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"payload": result})
}

func (h *ClassHandler) GetClassByID(ctx *gin.Context) {
	// Get id of class
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := h.repo.GetOne(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting class"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"payload": result})

}

func (h *ClassHandler) CreateClass(ctx *gin.Context) {

	// Verify Admin Role
	role, _ := ctx.Get("role")
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Only Admin can create classes"})
		return
	}

	// Get json input values
	var json class.CreateClassDTO
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

	newClass := class.Class{
		NameClass:   json.NameClass,
		Professor:   json.Professor,
		Description: json.Description,
		DateClass:   dateClass,
		Capacity:    json.Capacity,
	}

	result, err := h.repo.Insert(&newClass)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Class already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating class"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"payload": result})
}

func (h *ClassHandler) UpdateClass(ctx *gin.Context) {

	// Verify Admin Role
	role, _ := ctx.Get("role")
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Only Admin can update classes"})
		return
	}

	// Get id of class
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get json input values
	var json class.UpdateClassDTO
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error when binding json"})
		return
	}

	// Only update non-zero fields by default
	updatedValues := class.Class{
		NameClass:   json.NameClass,
		Professor:   json.Professor,
		Description: json.Description,
		Capacity:    json.Capacity,
	}

	// Parse DateTime if sent
	if dateClass, err := time.Parse(time.DateTime, json.DateClass); json.DateClass != "" {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date/time format. Use YYYYY-MM-DD HH:MM:SS"})
			return
		}
		updatedValues.DateClass = dateClass
	}

	result, err := h.repo.Update(id, &updatedValues)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating class"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"payload": result})
}

func (h *ClassHandler) DeleteClass(ctx *gin.Context) {

	// Verify Admin Role
	role, _ := ctx.Get("role")
	if role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Only Admin can delete classes"})
		return
	}

	// Get id of class
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting class"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"payload": "The Class has been deleted succesfully"})
}
