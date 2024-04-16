package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

// CreateComment creates a new comment for a post.
//
// This endpoint allows authenticated users to create a new comment for a specific post by providing the post ID and comment content.
// The comment is associated with the authenticated user who created it.
//
// Authentication:
//   - This endpoint requires the user to be authenticated. The user must include a valid JWT token in the request headers.
//
// Parameters:
//   - c: Context: The Gin context containing information about the HTTP request.
//
// Expected JSON Request Body:
//
//	{
//	  "content": "string" // The content of the comment.
//	}
//
// Responses:
//   - 201 Created: Returns the created comment object if successful.
//   - 400 Bad Request: If the request body is invalid or if the post ID is invalid.
//   - 401 Unauthorized: If the user is not authenticated or the JWT token is invalid or expired.
//   - 500 Internal Server Error: If an unexpected error occurs while creating the comment.
//
// @Summary Create a new comment
// @Description Create a new comment for a post
// @Tags Comments
// @Security JwtAuth
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Param comment body models.Comment true "Comment object"
// @Success 201 {object} models.Comment
// @Router /posts/{postId}/comments [post]
func CreateComment(c *gin.Context) {
	postIdStr := c.Param("postId")
	postIdUint, err := strconv.ParseUint(postIdStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid postId"})
		return
	}
	postId := uint(postIdUint)

	var newComment models.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
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

	newComment.PostID = postId

	// Check if the post exists
	var post models.Post
	if err := initializer.DB.First(&post, newComment.PostID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post_id"})
		return
	}

	newComment.UserID = userModel.ID

	if err := initializer.DB.Create(&newComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, newComment)
}

func GetCommentsForPost(c *gin.Context) {
	postId := c.Param("postId")
	var comments []models.Comment
	if err := initializer.DB.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentID")

	if err := initializer.DB.Where("id = ?", commentID).Delete(&models.Comment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func UpdateComment(c *gin.Context) {
	var updatedComment models.Comment
	commentID := c.Param("commentID")

	if err := initializer.DB.First(&updatedComment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := initializer.DB.Save(&updatedComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, updatedComment)
}
