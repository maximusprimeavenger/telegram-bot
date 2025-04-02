package repository

import (
	"telegram-bot/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RemoveNotifications(id int, email string) error {
	result := r.db.Model(&models.User{}).Where(&models.User{Id: id}).Update("notify_mode", false)
	return result.Error
}
