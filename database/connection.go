package database

import (
	"log"
	"os"

	"github.com/kenjius01/social-sever/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {

	db_url := os.Getenv("DB_URL")

	connection, err := gorm.Open(mysql.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatal("Fail to load database!\n", err.Error())
	}
	log.Println("Connect to the database successfully")
	connection.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migration!")
	// Add migration

	connection.AutoMigrate(&models.User{}, &models.Follower{}, &models.Post{}, &models.Like{})
	DB = connection

}
