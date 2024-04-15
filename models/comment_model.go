package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID  uint   `json:"post_id"`
	UserID  uint   `json:"user_id"` // Optional
	Content string `json:"content"`
	Post    Post   `gorm:"foreignKey:PostID;references:ID"` // Many-to-One relationship with Post
	User    User   `gorm:"foreignKey:UserID;references:ID"` // Optional, Many-to-One relationship with User
}
