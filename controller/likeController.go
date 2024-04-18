package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

// LikePost adds a like to a post.
// @Summary Like a post
// @Description Like a post
// @Tags Like
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param postId path uint true "Post ID"
// @Success 201 {string} string "Post liked successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{postId}/likes [POST]
func LikePost(c *gin.Context) {
	postIdStr := c.Param("postId")
	postIdUint, err := strconv.ParseUint(postIdStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid postId"})
		return
	}
	postId := uint(postIdUint)

	// Check if the user already liked the post
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

	var existingLike models.Like
	if err := initializer.DB.Where("post_id = ? AND user_id = ?", postId, userModel.ID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already liked the post"})
		return
	}

	// Add like to the post
	like := models.Like{
		PostID: postId,
		UserID: userModel.ID,
	}
	if err := initializer.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like the post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post liked successfully"})
}

// UnlikePost removes a like from a post.
// @Summary Unlike a post
// @Description Unlike a post
// @Tags Like
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param postId path uint true "Post ID"
// @Success 200 {string} string "Post unliked successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{postId}/likes [DELETE]
func UnlikePost(c *gin.Context) {
	postId := c.Param("postId")

	// Get the user from context
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

	// Delete the like
	if err := initializer.DB.Where("post_id = ? AND user_id = ?", postId, userModel.ID).Delete(&models.Like{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike the post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

// GetLikesForPost retrieves the number of likes for a post.
// @Summary Get likes for a post
// @Description Get likes for a post
// @Tags Like
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param postId path uint true "Post ID"
// @Success 200 {object} gin.H "Likes count"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{postId}/likes [GET]
func GetLikesForPost(c *gin.Context) {
	postId := c.Param("postId")
	var likesCount int64
	if err := initializer.DB.Model(&models.Like{}).Where("post_id = ?", postId).Count(&likesCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch likes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"likesCount": likesCount})
}
