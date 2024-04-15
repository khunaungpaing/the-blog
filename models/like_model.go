package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
	User   User `gorm:"foreignKey:UserID;references:ID"` // Many-to-One relationship with User
	Post   Post `gorm:"foreignKey:PostID;references:ID"` // Many-to-One relationship with Post
	// You can add additional fields like 'liked_at' for timestamp if needed
}
