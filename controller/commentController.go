package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
)

// CreateComment creates a new comment for a post.
// @Summary Create a new comment
// @Description Create a new comment for a post
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param newComment body models.Comment true "New Comment object"
// @Success 201 {object} models.Comment "Created comment"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /comments [POST]
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

// GetCommentsForPost retrieves comments for a post.
// @Summary Get comments for a post
// @Description Get comments for a post
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param postId path uint true "Post ID"
// @Success 200 {array} models.Comment "Comments"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /comments/{postId} [GET]
func GetCommentsForPost(c *gin.Context) {
	postId := c.Param("postId")
	var comments []models.Comment
	if err := initializer.DB.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// DeleteComment deletes a comment by ID.
// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param commentID path uint true "Comment ID"
// @Success 200 {string} string "Comment deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Comment not found"
// @Failure 500 {string} string "Internal server error"
// @Router /comments/{commentID} [DELETE]
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentID")

	if err := initializer.DB.Where("id = ?", commentID).Delete(&models.Comment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// UpdateComment updates a comment by ID.
// @Summary Update a comment
// @Description Update a comment by ID
// @Tags Comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Param commentID path uint true "Comment ID"
// @Param updatedComment body models.Comment true "Updated Comment object"
// @Success 200 {object} models.Comment "Updated comment"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Comment not found"
// @Failure 500 {string} string "Internal server error"
// @Router /comments/{commentID} [PUT]
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
