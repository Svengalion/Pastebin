package models

import (
	"time"
)

type Paste struct {
	Hash      string    `json:"hash" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
