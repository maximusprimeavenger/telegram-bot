package db

import (
	"fmt"
	"log"
	"os"
	"telegram-bot/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var loadPath = "/app/.env"

func ConnectToSQL() (*gorm.DB, error) {
	err := godotenv.Load(loadPath)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(mysql:%s)/notifier?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"), os.Getenv("PASSWORD"), os.Getenv("PORT_MYSQL"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error with migrating database: %v", err)
	}
	db.AutoMigrate(&models.User{})
	log.Println("Connected to MySQL!")
	return db, err
}
