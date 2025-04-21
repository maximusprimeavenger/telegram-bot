package repository

import (
	"log"
	"telegram-bot/internal/db"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"

	"github.com/badoux/checkmail"
	"gorm.io/gorm"
)

func (client *Client) Auth(update models.Update, chatID int, userSession *models.UserSession, step, text string) error {
	log.Printf("Auth called for chatId: %d, Step: %s", chatID, step)
	log.Println("Connecting to SQL")
	dbConn, err := db.ConnectToSQL()
	if err != nil {
		log.Fatal("Database connection failed:", err)
		return err
	}

	var username string
	if update.Message != nil {
		username = update.Message.Chat.Username
	} else if update.CallbackQuery != nil {
		username = update.CallbackQuery.From.Username
	}

	var foundUser models.User
	if step == "start" {
		result := dbConn.Where("username = ?", username).First(&foundUser)

		if result.Error == nil {
			client.SendMessage(chatID, "Welcome back, "+foundUser.Name+"!")
			return nil
		}

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			log.Println("Database error:", result.Error)
			return result.Error
		}
	}

	if _, exists := client.UserSessions[chatID]; !exists {
		client.UserSessions[chatID] = &models.UserSession{
			User: models.User{
				Username: username,
				ID:       chatID,
			},
			Step: "name",
		}

		return nil
	}

	switch step {
	case "start":
		client.SendMessageWithButtons(chatID, "Enter your name", []string{"Cancel"}, []string{"cancel"})
		userSession.Step = "name"
	case "name":
		if text != "" {
			userSession.User.Name = text
		}
		userSession.Step = "email"
		client.SendMessageWithButtons(chatID, "Enter your email", []string{"Cancel"}, []string{"cancel"})
	case "email":
		log.Println("Hey, I'm started email!")
		if !isValidEmail(text) || text == "" {
			client.SendMessage(chatID, "Email is not valid")
			return nil
		}

		userSession.Step = "done"
		userId, err := helpers.NotifierIdTaking(text)
		if err != nil {
			log.Print(err)
		}

		var user = models.User{
			ID:         chatID,
			Username:   username,
			Name:       userSession.User.Name,
			Email:      text,
			NotifierID: userId,
			NotifyMode: true,
		}
		log.Println("Hey, I'm here, under user!")
		result := dbConn.Create(&user)
		if result.Error != nil {
			log.Println("Error creating user:", result.Error)
		} else {
			log.Println("User created successfully: ", user)
		}
		client.SendMessage(chatID, "Registered successfully!")
		return result.Error
	}

	return nil
}

func isValidEmail(email string) bool {
	err := checkmail.ValidateFormat(email)
	return err == nil
}
