package models

import (
	"time"
)

// Model DB Class
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

// Model DB User
type User struct {
	ID        uint   `json:"id"`
	NameUser  string `json:"name"`
	Email     string `json:"email"`
	Password  string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Model when creating User
type CreateUserDTO struct {
	NameUser string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Model when updating User
type UpdateUserDTO struct {
	NameUser string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Model when login User
type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Model for response User
type ResponseUser struct {
	ID        uint      `json:"id"`
	NameUser  string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
