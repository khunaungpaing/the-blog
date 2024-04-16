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

// @title the blog api
// @description Your API Description
// @version 1.0
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	post := v1.Group("/posts")
	{
		post.POST("/", middleware.RequireAuth, controller.CreatePost)
		post.GET("/", middleware.RequireAuth, controller.GetPosts)
		post.GET("/:postId", middleware.RequireAuth, controller.GetPost)
		post.DELETE("/:postId", middleware.RequireAuth, controller.DeletePost)
		post.PUT("/:postId", middleware.RequireAuth, controller.UpdatePost)
	}

	cmt := post.Group("/:postId/comments")
	{
		cmt.POST("/", middleware.RequireAuth, controller.CreateComment)
		cmt.GET("/", middleware.RequireAuth, controller.GetCommentsForPost)
		cmt.DELETE("/:commentId", middleware.RequireAuth, controller.DeleteComment)
		cmt.PUT("/:commentId", middleware.RequireAuth, controller.UpdateComment)
	}

	like := post.Group("/:postId/likes")
	{
		like.POST("/", middleware.RequireAuth, controller.LikePost)
		like.GET("/", middleware.RequireAuth, controller.GetLikesForPost)
		like.DELETE("/", middleware.RequireAuth, controller.UnlikePost)
	}

	v1.POST("/signup", controller.SignUp)
	v1.POST("/login", controller.Login)

	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Run()
}
