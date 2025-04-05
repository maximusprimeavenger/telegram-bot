package repository

import (
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) findUser(id int) (*models.User, error) {
	var user = models.User{}
	result := r.db.Find(&user).Where(&models.User{ID: id})
	if result.Error != nil {
		return nil, helpers.ErrorHelper(result.Error, "Couldn't find user")
	}
	return &user, nil
}
