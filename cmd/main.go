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
			if update.Message == nil && update.CallbackQuery == nil {
				continue
			}
			id := getId(&update)
			session := getSession(id, client)
			go client.RunNotifications(id, session)
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

func getId(update *models.Update) int {
	if update.Message != nil {
		return update.Message.Chat.Id
	} else {
		return update.CallbackQuery.Message.Chat.Id
	}
}

func getSession(id int, client *repository.Client) *models.UserSession {
	session, exists := client.UserSessions[id]
	if !exists {
		session = &models.UserSession{Step: "start"}
		client.UserSessions[id] = session
	}
	return session
}
