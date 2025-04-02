package models

type UserSession struct {
	User User
	Step string
}

type User struct {
	Id         int    `gorm:"primaryKey"`
	Email      string `gorm:"uniqueIndex;not null"`
	Name       string `gorm:"size:30;not null"`
	Username   string `gorm:"uniqueIndex;not null"`
	NotifierId string `gorm:"uniqueIndex;not null"`
	NotifyMode bool
}
