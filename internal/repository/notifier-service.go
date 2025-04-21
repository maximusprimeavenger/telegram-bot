package repository

import (
	"log"
	"telegram-bot/internal/db"
	"telegram-bot/internal/models"
)

func NotificationsOnOff(id int, answer bool) error {
	dbConn, err := db.ConnectToSQL()
	if err != nil {
		return err
	}
	repo := NewUserRepository(dbConn)
	err = repo.RemoveTurnNotifications(id, answer)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) RemoveTurnNotifications(id int, mode bool) error {
	result := r.db.Model(&models.User{}).Where(&models.User{ID: id}).Update("notify_mode", mode)
	return result.Error
}

func NotifyMode(id int) (string, string) {
	dbConn, err := db.ConnectToSQL()
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	repo := NewUserRepository(dbConn)
	user, err := repo.findUser(id)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	if user.NotifyMode {
		return "on", "off"
	} else {
		return "off", "on"
	}
}
