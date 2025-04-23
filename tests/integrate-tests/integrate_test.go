package tests_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	loadPath string = "/app/.env"
)

type User struct {
	Id    int    `gorm:"primaryKey"`
	Email string `gorm:"type:varchar(50);not null"`
	Name  string `gorm:"type:varchar(15);not null"`
}

func TestConnectToSQL(t *testing.T) {
	err := godotenv.Load(loadPath)
	if err != nil {
		t.Fatalf("error loading environment variables: %v", err)
	}

	dsn := fmt.Sprintf("testuser:testpass@tcp(mysql:%s)/notifier_test?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("PORT_MYSQL"))
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("error with migrating database: %v", err)
	}

	err = dbConn.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("error during auto migration: %v", err)
	}

	t.Logf("Connected to MySQL and migrated successfully!")

	t.Run("Create user", func(t *testing.T) {
		user := &User{
			Name:  "Test",
			Email: "test@example.com",
		}
		result := dbConn.Create(&user)
		if result.Error != nil {
			t.Fatalf("error creating user: %v", result.Error)
		}

		assert.NotZero(t, user.Id, "User ID should be assigned")
		var foundUser User
		if err = dbConn.First(&foundUser, "email = ?", "test@example.com").Error; err != nil {
			t.Fatalf("error fetching user: %v", err)
		}
		assert.Equal(t, user.Email, foundUser.Email, "User's email does not match")
		assert.Equal(t, user.Name, foundUser.Name, "User's name does not match")
	})

	t.Run("Update user", func(t *testing.T) {
		var user, foundUser User
		if err := dbConn.First(&user, "email = ?", "test@example.com").Error; err != nil {
			t.Fatalf("error finding user: %v", err)
		}
		user.Name = "TEST"
		result := dbConn.Save(&user)
		if result.Error != nil {
			t.Fatalf("error updating user: %v", result.Error)
		}

		if err = dbConn.First(&foundUser, "email = ?", "test@example.com").Error; err != nil {
			t.Fatalf("error finding user: %v", err)
		}
		assert.Equal(t, "TEST", foundUser.Name, "Name has not been updated")
	})
	t.Run("Delete user", func(t *testing.T) {
		var user User
		if err := dbConn.First(&user, "email = ?", "test@example.com").Error; err != nil {
			t.Fatalf("error fetching user: %v", err)
		}
		result := dbConn.Delete(&user)
		if result.Error != nil {
			t.Fatalf("error deleting user: %v", result.Error)
		}
		if err := dbConn.First(&user, "email = ?", "test@example.com").Error; err == nil {
			t.Fatalf("expected error fetching deleted user, but found: %v", user)
		}
	})
	dbConn.Migrator().DropTable(&User{})
}
