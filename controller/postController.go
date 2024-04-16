package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

func CreatePost(c *gin.Context) {
	// 1. Bind incoming JSON data to a Post struct
	var newPost models.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. Get the user from the context
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found in context"})
		return
	}

	// Assert user to models.User
	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	// Assign UserID to the new post
	newPost.UserID = userModel.ID

	// 3. Save the Post struct
	if err := initializer.DB.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// 4. Return the created Post struct
	c.JSON(http.StatusCreated, newPost)
}

// @Summary Get a post by ID
// @Description Retrieve a post information based on its unique identifier. Requires authentication.
// @ID get-post-by-id
// @Tags posts
// @Accept  json
// @Produce  json
// @Param postId path int true "Post ID"
// @Success 200 {object} models.Post "ok"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Post not found"
// @Security apiKey [read] // Assuming your middleware uses an "apiKey" scheme with read scope
// @Router /posts/:postId [get]
func GetPost(c *gin.Context) {
	var post models.Post
	id := c.Param("postId")
	if err := initializer.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func GetPosts(c *gin.Context) {
	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated posts
	var posts []models.Post
	var totalPostsCount int64
	if err := initializer.DB.Model(&models.Post{}).Count(&totalPostsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total posts count"})
		return
	}
	if err := initializer.DB.Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// Prepare response with pagination metadata
	response := gin.H{
		"posts":       posts,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalCount":  totalPostsCount,
	}

	c.JSON(http.StatusOK, response)
}

func UpdatePost(c *gin.Context) {
	var updatedPost models.Post
	id := c.Param("postId")
	if err := initializer.DB.First(&updatedPost, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := initializer.DB.Save(&updatedPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, updatedPost)
}

func DeletePost(c *gin.Context) {
	var post models.Post
	id := c.Param("postId")
	if err := initializer.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	if err := initializer.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
