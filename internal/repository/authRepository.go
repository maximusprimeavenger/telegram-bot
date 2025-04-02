package repository

import (
	"log"
	"regexp"
	"telegram-bot/internal/db"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"

	"gorm.io/gorm"
)

func (client *Client) Auth(update models.Update, chatID int, user *models.UserSession, step, text string) error {
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
				Id:       chatID,
			},
			Step: "name",
		}

		return nil
	}

	session := client.UserSessions[chatID]

	switch step {
	case "start":
		client.SendMessageWithButtons(chatID, "Enter your name", []string{"Cancel"}, []string{"cancel"})
		session.Step = "name"
	case "name":
		if text != "" {
			session.User.Name = text
		}
		session.Step = "email"
		client.SendMessageWithButtons(chatID, "Enter your email", []string{"Cancel"}, []string{"cancel"})
	case "email":
		if !isValidEmail(session.User.Email) || text == "" {
			return helpers.ErrorHelper(err, "Email not valid")
		}
		session.Step = "done"
		userId, err := helpers.NotidierIdTaking(session.User.Email)
		if err != nil {
			return err
		}
		var user = models.User{
			Id:         chatID,
			Username:   username,
			Name:       session.User.Name,
			Email:      session.User.Email,
			NotifierId: userId,
			NotifyMode: true,
		}

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
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
