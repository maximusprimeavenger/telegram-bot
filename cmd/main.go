package main

import (
	"log"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"
)

const (
	host = "api.telegram.org"
)

func main() {
	token, err := helpers.Token()
	if err != nil {
		log.Fatal(err)
	}
	client := helpers.New(host, token)
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
			id, err := helpers.IdTaking(update)
			if err != nil {
				log.Fatalf("error with finding id: %v", err)
			}
			session, exists := client.UserSessions[id]
			if !exists {
				log.Println("Creating new session for chat:", id)
				session = &models.UserSession{Step: "start"}
				client.UserSessions[id] = session
			}

			// Button touching
			if update.CallbackQuery != nil {
				queryChatId := update.CallbackQuery.Message.Chat.Id
				data := update.CallbackQuery.Data

				if data == "" {
					log.Printf("Empty callback query data from chat: %d", queryChatId)
					continue
				}

				switch data {
				case "help":
					client.SendMessage(queryChatId, "Available commands:\n/my_orders - View orders\n/register - Register an account")

				case "check_orders":
					client.SendMessage(queryChatId, "Sorry, for now it's not available")

				case "answer_yes":
					client.SendMessage(id, "Ok, now, you need to authorize.")
					log.Printf("Session for chat %d: %+v", queryChatId, session)
					client.Auth(update, id, session, "start", "")

				case "answer_no":
					client.SendMessage(id, "Ok, you can authorize later.")
				case "cancel":
					client.SendMessage(id, "Canceled")

				default:
					log.Printf("Unknown callback data: %s", data)
				}

			}

			// Text messages
			if update.Message != nil {
				messageChatId := update.Message.Chat.Id
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
					client.SendMessage(messageChatId, "Your list of users:")
				case "Check my orders":
					client.SendMessage(messageChatId, "Your list of orders:")
				case "Help":
					client.SendMessage(messageChatId, "Available commands:\n/my_orders - View orders\n/register - Register an account")
				default:
					switch session.Step {
					case "name":
						client.Auth(update, id, session, "name", update.Message.Text)
					case "email":
						client.Auth(update, id, session, "email", update.Message.Text)
					default:
						client.SendMessageWithButtons(messageChatId, "Please, follow these buttons\n 👇 👇 👇 👇",
							[]string{"Help", "Cancel"},
							[]string{"help", "cancel"},
						)
					}
				}
			}
		}
	}
}
