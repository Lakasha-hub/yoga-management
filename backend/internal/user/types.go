package user

import "time"

type (
	User struct {
		ID        uint      `json:"id"`
		NameUser  string    `json:"name"`
		Email     string    `json:"email"`
		Password  string    `json:"-" gorm:"column:password"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	RegisterUserDTO struct {
		NameUser string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	UpdateUserDTO struct {
		NameUser string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	LoginUserDTO struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
