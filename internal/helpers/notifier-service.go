package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"telegram-bot/internal/models"
)

func SendToNotifier() (*models.Order, error) {
	id, err := sendToUserAuth()
	if err != nil {
		return nil, ErrorHelper(err, "error sending request to user-auth-service")
	}
	resp, err := http.Post("http://notifier-service:8082/telegram-bot", "application/json", strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, id)))
	if err != nil {
		return nil, ErrorHelper(err, "error sending request to notifier-service")
	}

	defer resp.Body.Close()
	var order *models.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, ErrorHelper(err, "error parsing orders from JSON")
	}
	return order, nil
}

func sendToUserAuth() (string, error) {

	//resp, err := http.Post("http://user-auth-service:8081/telegram-bot/user", "application/json")
	return "", nil
}
