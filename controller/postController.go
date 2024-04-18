package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

// CreatePost creates a new post.
// @Summary Create a new post
// @Description Create a new post
// @Tags Post
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param newPost body models.Post true "New Post object"
// @Success 201 {object} models.Post "Created post"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /posts [POST]
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

// GetPost retrieves a single post.
//
// @Summary Get a single post
// @Description Get a single post
// @Tags Post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Success 200 {object} models.Post "Post"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{postId} [GET]
func GetPost(c *gin.Context) {
	var post models.Post
	id := c.Param("postId")
	if err := initializer.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// GetPosts retrieves paginated posts.
//
// @Summary Get paginated posts
// @Description Get paginated posts
// @Tags Post
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} models.Post "Paginated posts"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /posts [GET]
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

// UpdatePost updates an existing post.
// @Summary Update an existing post
// @Description Update an existing post
// @Tags Post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param updatedPost body models.Post true "Updated Post object"
// @Success 200 {object} models.Post "Updated post"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Post not found"
// @Failure 403 {string} string "User is not authorized to update this post"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{postId} [PUT]
func UpdatePost(c *gin.Context) {
	// Get the post ID from the URL parameter
	id := c.Param("postId")

	// Bind the incoming JSON data to the Post struct
	var updatedPost models.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the post exists
	var post models.Post
	if err := initializer.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if the user is authorized to update the post
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found in context"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	if post.UserID != userModel.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not authorized to update this post"})
		return
	}

	// Update only the fields that are allowed to be updated
	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	post.Status = updatedPost.Status

	// Save the updated post
	if err := initializer.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Return the updated post
	c.JSON(http.StatusOK, post)
}

// DeletePost deletes an existing post.
// @Summary Delete an existing post
// @Description Delete an existing post
// @Tags Post
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Success 200 {object} gin.H "Post deleted successfully"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 403 {object} gin.H "User is not authorized to delete this post"
// @Failure 500 {object} gin.H "Failed to delete post"
// @Router /posts/{postId} [DELETE]
func DeletePost(c *gin.Context) {
	var post models.Post
	id := c.Param("postId")
	if err := initializer.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if the user is authorized to delete the post
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found in context"})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	if post.UserID != userModel.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not authorized to delete this post"})
		return
	}

	if err := initializer.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
