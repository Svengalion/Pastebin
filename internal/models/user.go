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
