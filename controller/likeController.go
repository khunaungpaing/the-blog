package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

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

func GetLikesForPost(c *gin.Context) {
	postId := c.Param("postId")
	var likesCount int64
	if err := initializer.DB.Model(&models.Like{}).Where("post_id = ?", postId).Count(&likesCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch likes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"likesCount": likesCount})
}
