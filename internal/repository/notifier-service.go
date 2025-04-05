package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"telegram-bot/internal/db"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"
)

func SendToNotifier(id string) (*models.Order, error) {
	resp, err := http.Post("http://notifier-service:8082/telegram-bot", "application/json", strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, id)))
	if err != nil {
		return nil, helpers.ErrorHelper(err, "error sending request to notifier-service")
	}

	defer resp.Body.Close()
	var order *models.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, helpers.ErrorHelper(err, "error parsing orders from JSON")
	}
	return order, nil
}

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
