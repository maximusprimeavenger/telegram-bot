package handlers

import (
	"log"
	"telegram-bot/internal/models"
	"telegram-bot/internal/repository"
)

func HandleCallback(update models.Update, client *repository.Client) {
	queryChatId := update.CallbackQuery.Message.Chat.Id
	session := getSession(queryChatId, client)

	data := update.CallbackQuery.Data

	if data == "" {
		log.Printf("Empty callback query data from chat: %d", queryChatId)
	}

	switch data {
	case "help":
		client.SendMessage(queryChatId, "Available commands:\n/my_orders - View orders\n/register - Register an account")

	case "check_orders":
		client.SendMessage(queryChatId, "Sorry, for now it's not available")

	case "answer_yes":
		client.SendMessage(queryChatId, "Ok, now, you need to authorize.")
		log.Printf("Session for chat %d: %+v", queryChatId, session)
		client.Auth(update, queryChatId, session, "start", "")

	case "answer_no":
		client.SendMessage(queryChatId, "Ok, you can authorize later.")
	case "cancel":
		delete(client.UserSessions, queryChatId)
		client.SendMessage(queryChatId, "Canceled")

	default:
		log.Printf("Unknown callback data: %s", data)
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
