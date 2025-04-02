package handlers

import (
	"log"
	"telegram-bot/internal/models"
	"telegram-bot/internal/repository"
)

func HandleMessage(update models.Update, client *repository.Client) {
	messageChatId := update.Message.Chat.Id
	session := getSession(messageChatId, client)
	log.Printf("New message from chat %d: %s", messageChatId, update.Message.Text)

	switch update.Message.Text {
	case "/start":
		client.SendMessageWithButtons(messageChatId,
			"Hello! Welcome to service of notifications. Want to authorize to receive notifications?",
			[]string{"Yes", "No"},
			[]string{"answer_yes", "answer_no"},
		)
	case "/my_orders":
		client.SendMessage(messageChatId, "Your list of orders:")
	case "/register":

	case "Check my orders":
		client.SendMessage(messageChatId, "Your list of orders:")
	case "Help":
		client.SendMessage(messageChatId, "Available commands:\n/my_orders - View orders\n/register - Register an account")
	default:
		switch session.Step {
		case "name":
			client.Auth(update, messageChatId, session, "name", update.Message.Text)
		case "email":
			client.Auth(update, messageChatId, session, "email", update.Message.Text)
		default:
			client.SendMessageWithButtons(messageChatId, "Please, follow these buttons\n 👇 👇 👇 👇",
				[]string{"Help", "Cancel"},
				[]string{"help", "cancel"},
			)
		}
	}
}
