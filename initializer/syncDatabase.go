package initializer

import "github.com/khunaungpaing/the-blog-api/models"

func SyncDataBase() {
	DB.AutoMigrate(&models.User{})
}
