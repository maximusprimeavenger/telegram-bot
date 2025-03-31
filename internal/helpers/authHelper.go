package helpers

import (
	"log"
	"regexp"
	"telegram-bot/internal/db"
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

	log.Println("Database connected successfully!")
	var username string
	if update.Message != nil {
		username = update.Message.Chat.Username
	} else if update.CallbackQuery != nil {
		username = update.CallbackQuery.From.Username
	}

	log.Printf("Calling Auth for chat: %d, user:%v", chatID, user)
	var foundUser models.User
	result := dbConn.Where("username = ?", username).First(&foundUser)

	if result.Error == nil {
		client.UserSessions[chatID] = &models.UserSession{
			User: models.User{
				Id:       chatID,
				Username: foundUser.Username,
				Name:     foundUser.Name,
				Email:    foundUser.Email,
			},
			Step: "done",
		}
		client.SendMessage(chatID, "Welcome back, "+foundUser.Name+"!")
		return nil
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Println("Database error:", result.Error)
		return result.Error
	}

	if _, exists := client.UserSessions[chatID]; !exists {
		client.UserSessions[chatID] = &models.UserSession{
			User: models.User{
				Username: update.Message.Chat.Username,
				Id:       chatID,
			},
			Step: "name",
		}

		return nil
	}

	session := client.UserSessions[chatID]

	switch step {
	case "start":
		session.User.Id = chatID
		client.SendMessageWithButtons(chatID, "Enter your name", []string{"Cancel"}, []string{"cancel"})
		session.Step = "name"
	case "name":
		if text != "" {
			session.User.Name = text
		}
		session.Step = "email"
		client.SendMessageWithButtons(chatID, "Enter your email", []string{"Cancel"}, []string{"cancel"})
	case "email":
		if text != "" {
			session.User.Email = text
		}
		if !isValidEmail(session.User.Email) {
			return ErrorHelper(err, "Email not valid")
		}
		session.Step = "done"
		var user = models.User{
			Id:       chatID,
			Username: username,
			Name:     session.User.Name,
			Email:    session.User.Email,
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
