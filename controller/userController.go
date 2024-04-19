package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/khunaungpaing/the-blog-api/dto"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/models"
	"golang.org/x/crypto/bcrypt"
)

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param body body models.User true "User information"
// @Success 201 {object} models.User
// @Failure 400 {object} string
// @Router /signup [post]
func SignUp(c *gin.Context) {
	// Get the email/password from the request body
	var body models.RegisterRequest

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	// Create a new user
	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}

	// Save the user
	result := initializer.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	user.Password = ""
	// Return the user
	c.JSON(http.StatusCreated, user)

}

// Login godoc
// @Summary Login an existing user
// @Description Login an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "User login information"
// @Success 200 {object} models.User
// @Failure 401 {object} string
// @Router /login [post]
func Login(c *gin.Context) {
	// Get the email/password from the request body
	var body models.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// look up the requested user
	var user models.User
	if err := initializer.DB.First(&user, "email =?", body.Email).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve user",
		})
		return
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Cannot create token",
		})
		return
	}

	// send the token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

// GetUserProfile returns the profile of the logged-in user.
// @Summary Get user profile
// @Description Returns the profile of the logged-in user.
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Success 200 {object} models.User "User profile"
// @Failure 400 {object} gin.H "User not found in context"
// @Failure 500 {object} gin.H "Failed to get user from context"
// @Router /users/profile [get]
func GetUserProfile(c *gin.Context) {
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

	// Response user
	userModel.Password = "" // Ensure password is not sent in response
	c.JSON(http.StatusOK, userModel)
}

// UpdateUserProfile updates the profile of the logged-in user.
// @Summary Update user profile
// @Description Updates the profile of the logged-in user.
// @Tags Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token using the Bearer scheme"
// @Param body body dto.RequestUser true "User object"
// @Success 200 {object} models.User "Updated user profile"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 401 {object} gin.H "User not found in context"
// @Failure 500 {object} gin.H "Failed to get user from context"
// @Router /users/profile [patch]
func UpdateUserProfile(c *gin.Context) {
	var body dto.RequestUser
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Update user profile
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

	// Update user profile fields
	userModel.Username = body.Username
	userModel.Email = body.Email
	userModel.Bio = body.Bio
	userModel.ProfilePic = body.ProfilePic

	// Hash the password if provided
	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to hash password"})
			return
		}
		userModel.Password = string(hash)
	}

	// Save the updated user profile
	result := initializer.DB.Save(&userModel)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	userModel.Password = "" // Ensure password is not sent in response
	c.JSON(http.StatusOK, userModel)
}
