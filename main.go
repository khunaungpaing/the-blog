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
// @description The Blog API is a robust and efficient solution for managing a blog platform. Developed using Golang and Gin, it embodies modern practices and technologies for seamless performance and scalability. Leveraging PostgreSQL as its database engine ensures reliability and flexibility in data management, while JWT (JSON Web Tokens) authentication enhances security by providing a stateless authentication mechanism. \n\nThis API serves as a foundational component for building and managing a dynamic blogging platform, offering a comprehensive set of endpoints for user authentication, post management, comment handling, user profile management, and more. With clear and concise documentation and a user-friendly architecture, integrating this API into your project is straightforward and hassle-free.\n\nWhether you're developing a personal blog, a collaborative writing platform, or an enterprise-level content management system, The Blog API provides the necessary tools and functionality to streamline your development process and deliver a seamless user experience. Unlock the power of modern web development with The Blog API.

// @contact.name API Support
// @contact.email khunaungpaing.it.tumlm@gmail.com

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
		post.GET("/", controller.GetPosts)
		// Get a specific post
		post.GET("/:postId", controller.GetPost)
		// Delete a specific post
		post.DELETE("/:postId", middleware.RequireAuth, controller.DeletePost)
		// Update a specific post
		post.PUT("/:postId", middleware.RequireAuth, controller.UpdatePost)
	}

	// Initialize the comments endpoint for a specific post
	postIdRoute := post.Group("/:postId")
	{
		// Create a new comment for a specific post
		postIdRoute.POST("/comments", middleware.RequireAuth, controller.CreateComment)
		// Get all the comments for a specific post
		postIdRoute.GET("/comments", middleware.RequireAuth, controller.GetCommentsForPost)
		// Delete a specific comment for a specific post
		postIdRoute.DELETE("/comments/:commentId", middleware.RequireAuth, controller.DeleteComment)
		// Update a specific comment for a specific post
		postIdRoute.PUT("/comments/:commentId", middleware.RequireAuth, controller.UpdateComment)

		// Like the comments endpoint
		postIdRoute.POST("/likes", middleware.RequireAuth, controller.LikePost)
		// Get all the likes for a specific post
		postIdRoute.GET("/likes", middleware.RequireAuth, controller.GetLikesForPost)
		// Unlike a specific post
		postIdRoute.DELETE("/likes", middleware.RequireAuth, controller.UnlikePost)
	}

	v1.POST("/users/signup", controller.SignUp)
	v1.POST("/users/login", controller.Login)
	v1.GET("/users/profile", middleware.RequireAuth, controller.GetUserProfile)
	v1.PUT("/users", middleware.RequireAuth, controller.UpdateUserProfile)

	// Set the trusted proxies
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Run()
}
