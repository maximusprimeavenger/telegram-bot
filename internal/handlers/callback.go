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
		client.SendMessage(queryChatId, "Available commands:\n/my_orders - View orders\n/register - Register an account\n/set_notify - Set notifications")

	case "check_orders":
		client.CheckOrders(update, queryChatId)

	case "answer_yes":
		client.Auth(update, queryChatId, session, "start", "")

	case "answer_no":
		client.SendMessage(queryChatId, "Ok, you can authorize later.")

	case "notify_no":
		repository.NotificationsOnOff(queryChatId, false)
		client.SendMessage(queryChatId, "Notifications are disabled")

	case "notify_yes":
		repository.NotificationsOnOff(queryChatId, true)
		client.SendMessage(queryChatId, "Notifications are enabled")

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
