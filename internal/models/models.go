package models

import (
	"time"
)

// Model DB
type Class struct {
	ID        uint      `json:"id"`
	NameClass string    `json:"name"`
	Professor string    `json:"professor"`
	DateClass time.Time `json:"date"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Model when creating a Class
type CreateClassDTO struct {
	NameClass string `json:"name" binding:"required"`
	Professor string `json:"professor" binding:"required"`
	DateClass string `json:"date" binding:"required"`
	Capacity  int    `json:"capacity"`
}

// Model when updating a Class
type UpdateClassDTO struct {
	NameClass string `json:"name"`
	Professor string `json:"professor"`
	DateClass string `json:"date"`
	Capacity  int    `json:"capacity"`
}
