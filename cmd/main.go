package main

import (
	"log"
	"telegram-bot/internal/handlers"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"
	"telegram-bot/internal/repository"
)

const (
	host = "api.telegram.org"
)

func main() {
	token, err := helpers.Token()
	if err != nil {
		log.Fatal(err)
	}
	client := repository.New(host, token)
	offset := 0
	for {
		updates, err := client.Updates(offset, 25)
		if err != nil {
			log.Println("Error taking notifications:", err)
			continue
		}

		for _, update := range updates {
			offset = update.Id + 1
			if client.UserSessions == nil {
				client.UserSessions = make(map[int]*models.UserSession)
			}
			if update.CallbackQuery != nil {
				handlers.HandleCallback(update, client)
			}
			if update.Message != nil {
				handlers.HandleMessage(update, client)
			}
		}
	}
}
