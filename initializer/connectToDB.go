package initializer

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectToDB connects to the database.
func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to db")
	}

	fmt.Printf("Success connect to DB")
}
