package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
