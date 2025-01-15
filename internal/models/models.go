package models

import (
	"time"
)

type Class struct {
	ID        uint      `json:"id"`
	NameClass string    `json:"name"`
	Professor string    `json:"professor"`
	DateClass time.Time `json:"date"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
