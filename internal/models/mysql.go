package models

type UserSession struct {
	Id       int    `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Name     string `gorm:"size:30;not null"`
	Username string `gorm:"unique;not null"`
	Step     string
}
