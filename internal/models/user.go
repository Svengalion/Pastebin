package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Login     string    `gorm:"unique;not null" json:"login" validate:"required,min=3,max=32"`
	Email     string    `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required, min=3, max=32"`
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=6, max=64"`
}

type RegisterResponce struct {
	Id        uint      `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
