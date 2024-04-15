package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

func CreatePost(c *gin.Context) {
	// 1. Bind incoming JSON data to a Post struct
	var newPost *models.Post

	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(newPost)
	fmt.Println(c.Get("user"))

	// 2. Create a new Post struct
	user, exist := c.Get("user")

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	newPost.User = user.(models.User)

	// 3. Save the Post struct
	if err := initializer.DB.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Return the Post struct
	c.JSON(http.StatusCreated, newPost)
}
