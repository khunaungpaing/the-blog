package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/dto"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

// CreateComment creates a new comment for a specific post.
// @Summary Create a new comment
// @Description Create a new comment for the specified post.
// @Tags comments
// @Accept json
// @Produce json
// @Param postId path int true "ID of the post to which the comment will be added"
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Param comment body dto.RequestComment true "Comment object containing content"
// @Success 201 {object} models.Comment "Successfully created comment"
// @Failure 400 {object} gin.H "Bad request, invalid postId or request body"
// @Failure 401 {object} gin.H "Unauthorized access, missing or invalid token"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts/{postId}/comments [post]
func CreateComment(c *gin.Context) {
	postIdStr := c.Param("postId")
	postIdUint, err := strconv.ParseUint(postIdStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid postId"})
		return
	}
	postId := uint(postIdUint)

	var requestCmt dto.RequestComment
	if err := c.ShouldBindJSON(&requestCmt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

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

	var newComment models.Comment
	newComment.PostID = postId
	newComment.Content = requestCmt.Content
	newComment.UserID = userModel.ID

	var post models.Post
	if err := initializer.DB.First(&post, newComment.PostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	if err := initializer.DB.Create(&newComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, newComment)
}

// GetCommentsForPost retrieves comments for a specific post.
// @Summary Get comments for a post
// @Description Retrieves comments for the specified post.
// @Tags comments
// @Accept json
// @Produce json
// @Param postId path int true "ID of the post"
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Success 200 {array} models.Comment "Successfully retrieved comments"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts/{postId}/comments [get]
func GetCommentsForPost(c *gin.Context) {
	postId := c.Param("postId")
	var comments []models.Comment
	if err := initializer.DB.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// DeleteComment deletes a specific comment.
// @Summary Delete a comment
// @Description Deletes the specified comment.
// @Tags comments
// @Accept json
// @Produce json
// @Param commentId path int true "ID of the comment"
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Success 200 {object} gin.H "Successfully deleted comment"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts/{postId}/comments/{commentId} [delete]
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	if err := initializer.DB.Where("id = ?", commentID).Delete(&models.Comment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// UpdateComment updates a specific comment.
// @Summary Update a comment
// @Description Updates the specified comment.
// @Tags comments
// @Accept json
// @Produce json
// @Param postId path int true "ID of the post"
// @Param commentId path int true "ID of the comment"
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Param comment body dto.RequestComment true "Updated comment object"
// @Success 200 {object} models.Comment "Successfully updated comment"
// @Failure 400 {object} gin.H "Bad request, invalid request body"
// @Failure 404 {object} gin.H "Comment not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /posts/{postId}/comments/{commentId} [put]
func UpdateComment(c *gin.Context) {
	var requestCmt dto.RequestComment
	var updatedComment models.Comment
	commentID := c.Param("commentId")

	if err := initializer.DB.First(&updatedComment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := c.ShouldBindJSON(&requestCmt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedComment.Content = requestCmt.Content

	if err := initializer.DB.Save(&updatedComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, updatedComment)
}
