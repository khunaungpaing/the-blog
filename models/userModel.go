package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string `json:"username" gorm:"unique"` // Unique username
	Email      string `json:"email" gorm:"unique"`    // Unique email
	Password   string `json:"password,omitempty"`     // Hashed password, excluded from JSON
	Bio        string `json:"bio"`                    // Optional user bio
	ProfilePic string `json:"profile_pic"`            // Optional profile picture URL
	Posts      []Post `gorm:"foreignKey:UserID"`      // One-to-Many relationship with Post
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
