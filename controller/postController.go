package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/dto"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

// CreatePost creates a new post.
// @Summary Create a new post
// @Description Create a new post with the provided data.
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Param post body dto.RequestPost true "Post data"
// @Success 201 {object} models.Post "Successfully created post"
// @Failure 400 {object} gin.H "Bad request, invalid request body"
// @Failure 401 {object} gin.H "Unauthorized access, missing or invalid token"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts [post]
func CreatePost(c *gin.Context) {
	// 1. Bind incoming JSON data to a Post struct
	var requestPost dto.RequestPost
	if err := c.ShouldBindJSON(&requestPost); err != nil {
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

	// 3. Assign UserID to the new post
	var newPost models.Post
	newPost = requestPost.ToModel(newPost)
	newPost.UserID = userModel.ID
	newPost.User = nil

	// 4. Save the Post struct
	if err := initializer.DB.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post", "details": err.Error()})
		return
	}

	// 5. Return the created Post struct
	c.JSON(http.StatusCreated, newPost)
}

// GetPost retrieves a specific post by ID.
// @Summary Retrieve a specific post
// @Description Retrieve a specific post by its ID.
// @Tags posts
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Success 200 {object} models.Post "Found post"
// @Failure 404 {object} gin.H "Post not found"
// @Router /posts/{postId} [get]
func GetPost(c *gin.Context) {
	var post models.Post
	id := c.Param("postId")
	if err := initializer.DB.Preload("Categories").Preload("Tags").Preload("Comments").Preload("Media").Preload("User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	post.User.Password = ""
	c.JSON(http.StatusOK, post)
}

// GetPosts retrieves a list of posts with pagination.
// If the author query parameter is provided, it filters posts by that author.
// @Summary Retrieve a list of posts
// @Description Retrieve a list of posts with pagination support.
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Number of items per page (default: 10)"
// @Param author query string false "Filter posts by author username"
// @Success 200 {object} gin.H "List of posts"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts [get]
func GetPosts(c *gin.Context) {
	// Parse query parameters for pagination and author filtering
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	author := c.Query("author")

	// Calculate offset
	offset := (page - 1) * pageSize

	// Prepare DB query with or without author filtering
	dbQuery := initializer.DB.Preload("Categories").Preload("Tags").Preload("Comments").Preload("Media").Preload("User").Offset(offset).Limit(pageSize)
	if author != "" {
		// If author is provided, filter posts by author
		dbQuery = dbQuery.Where("user_id IN (SELECT id FROM users WHERE username = ?)", author)
	}

	// Fetch paginated posts with preloaded associations
	var posts []models.Post
	var totalPostsCount int64
	if err := dbQuery.Model(&models.Post{}).Count(&totalPostsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total posts count"})
		return
	}
	if err := dbQuery.Find(&posts).Error; err != nil {
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
// @Summary Update a post
// @Description Update an existing post with the provided data.
// @Tags posts
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Param post body dto.RequestPost true "Updated post data"
// @Success 200 {object} gin.H "Post updated successfully"
// @Failure 400 {object} gin.H "Bad request, invalid request body"
// @Failure 401 {object} gin.H "Unauthorized access, missing or invalid token"
// @Failure 403 {object} gin.H "Forbidden, user is not authorized to update this post"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts/{postId} [patch]
func UpdatePost(c *gin.Context) {
	// Get the post ID from the URL parameter
	id := c.Param("postId")

	// Bind the incoming JSON data to the Post struct
	var requestPost dto.RequestPost
	if err := c.ShouldBindJSON(&requestPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the post exists
	var post models.Post
	if err := initializer.DB.Preload("Categories").Preload("Tags").Preload("Comments").Preload("Media").First(&post, id).Error; err != nil {
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
	post = requestPost.ToModel(post)

	// Save the updated post
	if err := initializer.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated post
	c.JSON(http.StatusOK, post)
}

// DeletePost deletes a specific post by ID.
// @Summary Delete a post
// @Description Delete a specific post by its ID.
// @Tags posts
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Success 200 {object} gin.H "Post deleted successfully"
// @Failure 401 {object} gin.H "Unauthorized access, missing or invalid token"
// @Failure 403 {object} gin.H "Forbidden, user is not authorized to delete this post"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts/{postId} [delete]
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
