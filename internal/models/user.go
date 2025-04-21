package models

import "context"

type UserSession struct {
	User   User
	Step   string
	Cancel context.CancelFunc
}

type User struct {
	ID         int    `gorm:"primaryKey"`
	Email      string `gorm:"type:varchar(191);unique;not null"`
	Name       string `gorm:"type:varchar(30);not null"`
	Username   string `gorm:"type:varchar(191);unique;not null"`
	NotifierID string `gorm:"type:varchar(191);unique;not null"`
	NotifyMode bool   `gorm:"not null"`
}
