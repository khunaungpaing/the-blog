package main

import (
	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/controller"
	"github.com/khunaungpaing/the-blog-api/docs"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializer.LoadEnvVariable()
	initializer.ConnectToDB()
	initializer.SyncDataBase()
}

// main is the entry point of the application
// @title The Blog API
// @version 1.0
// @description This is a sample blog API
// @termsOfService https://example.com/terms/

// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email <EMAIL>

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	r := gin.Default()

	// Initialize the swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize the version 1 of the API
	v1 := r.Group("/api/v1")

	// Initialize the posts endpoint
	post := v1.Group("/posts")
	{
		// Create a new post
		post.POST("/", middleware.RequireAuth, controller.CreatePost)
		// Get all the posts
		post.GET("/", middleware.RequireAuth, controller.GetPosts)
		// Get a specific post
		post.GET("/:postId", middleware.RequireAuth, controller.GetPost)
		// Delete a specific post
		post.DELETE("/:postId", middleware.RequireAuth, controller.DeletePost)
		// Update a specific post
		post.PUT("/:postId", middleware.RequireAuth, controller.UpdatePost)
	}

	// Initialize the comments endpoint for a specific post
	cmt := post.Group("/:postId/comments")
	{
		// Create a new comment for a specific post
		cmt.POST("/", middleware.RequireAuth, controller.CreateComment)
		// Get all the comments for a specific post
		cmt.GET("/", middleware.RequireAuth, controller.GetCommentsForPost)
		// Delete a specific comment for a specific post
		cmt.DELETE("/:commentId", middleware.RequireAuth, controller.DeleteComment)
		// Update a specific comment for a specific post
		cmt.PUT("/:commentId", middleware.RequireAuth, controller.UpdateComment)
	}

	// Initialize the likes endpoint for a specific post
	like := post.Group("/:postId/likes")
	{
		// Like a specific post
		like.POST("/", middleware.RequireAuth, controller.LikePost)
		// Get all the likes for a specific post
		like.GET("/", middleware.RequireAuth, controller.GetLikesForPost)
		// Unlike a specific post
		like.DELETE("/", middleware.RequireAuth, controller.UnlikePost)
	}

	// Initialize the sign up endpoint
	v1.POST("/signup", controller.SignUp)
	// Initialize the sign in endpoint
	v1.POST("/login", controller.Login)

	// Set the trusted proxies
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Run()
}
