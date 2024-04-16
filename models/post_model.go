package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID      uint       `json:"user_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`            // Can include HTML for formatting
	Slug        string     `json:"slug" gorm:"unique"` // Unique URL slug for SEO
	Status      string     `json:"status" gorm:"type:string;check:status IN ('draft', 'published', 'archived')"`
	PublishedAt *time.Time `json:"published_at"`              // Optional publish timestamp
	UpdatedAt   *time.Time `json:"updated_at"`                // Timestamp of last update
	Categories  []Category `gorm:"many2many:post_categories"` // Optional, Many-to-Many relationship with Category (using a join table)
	Tags        []Tag      `gorm:"many2many:post_tags"`       // Optional, Many-to-Many relationship with Tag (using a join table)
	Comments    []Comment  `gorm:"foreignKey:PostID"`         // One-to-Many relationship with Comment
	Media       *Media     `gorm:"foreignKey:PostID"`         // Optional, One-to-One or One-to-Many relationship with Media
}

type Category struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `gorm:"many2many:post_categories"` // Optional, Many-to-Many relationship with Post (using a join table)
}

type Tag struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Posts       []Post `gorm:"many2many:post_tags"` // Optional, Many-to-Many relationship with Post (using a join table)
}

type Media struct {
	gorm.Model
	Filename string `json:"filename"`
	Path     string `json:"path"`      // Path to the media file in cloud storage
	MimeType string `json:"mime_type"` // Media file type (e.g., image/jpeg)
	PostID   uint   `json:"post_id"`   // Optional, Foreign Key referencing the Post
}
