package models

import "gorm.io/gorm"

// Comment represents a comment entity.
type Comment struct {
	gorm.Model
	PostID  uint   `json:"post_id"`
	UserID  uint   `json:"user_id"` // Optional
	Content string `json:"content"`
}
