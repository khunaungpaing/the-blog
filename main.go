package main

import (
	"github.com/gin-gonic/gin"
	"github.com/khunaungpaing/the-blog-api/controller"
	"github.com/khunaungpaing/the-blog-api/initializer"
	"github.com/khunaungpaing/the-blog-api/middleware"
)

func init() {
	initializer.LoadEnvVariable()
	initializer.ConnectToDB()
	initializer.SyncDataBase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controller.SignUp)
	r.POST("/login", controller.Login)
	r.POST("/post/create", middleware.RequireAuth, controller.CreatePost)
	r.GET("/home", middleware.RequireAuth, controller.Validate)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
