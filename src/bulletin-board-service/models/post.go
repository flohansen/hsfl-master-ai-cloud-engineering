package models

import (
	"time"
)

// Post struct represents a bulletin board post
type Post struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Title   string `json:"title"`
	Content string `json:"content"`
}
