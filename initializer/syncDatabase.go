package initializer

import "github.com/khunaungpaing/the-blog-api/models"

func SyncDataBase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.Tag{},
		&models.Media{},
		&models.Category{},
	)
}
