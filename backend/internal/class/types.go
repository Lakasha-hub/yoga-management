package class

import (
	"time"
)

type (
	Class struct {
		ID          uint      `json:"id"`
		NameClass   string    `json:"name"`
		Professor   string    `json:"professor"`
		DateClass   time.Time `json:"date"`
		Capacity    int       `json:"capacity"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
	CreateClassDTO struct {
		NameClass   string `json:"name" binding:"required"`
		Professor   string `json:"professor" binding:"required"`
		DateClass   string `json:"date" binding:"required"`
		Description string `json:"description" binding:"required"`
		Capacity    int    `json:"capacity"`
	}
	UpdateClassDTO struct {
		NameClass   string `json:"name"`
		Professor   string `json:"professor"`
		DateClass   string `json:"date"`
		Description string `json:"description"`
		Capacity    int    `json:"capacity"`
	}
)
