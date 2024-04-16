package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
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

	// Return the user
	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})

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
