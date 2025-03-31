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
	if client.UserSessions == nil {
		log.Println("UserSessions is nil, initializing...")
		client.UserSessions = make(map[int]*models.UserSession)
		user = &models.UserSession{
			Step: "start",
		}
	}
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
	// if update.Message != nil {
	// 	chatId = update.Message.Chat.Id
	// 	text = update.Message.Text
	// } else if update.CallbackQuery != nil {
	// 	chatId = update.CallbackQuery.Message.Chat.Id
	// } else {
	// 	log.Println("Error: update contains neither Message nor CallbackQuery")
	// 	return errors.New("invalid update structure")
	// }

	log.Printf("Calling Auth for chat: %d, user:%v", chatID, user)
	var foundUser models.UserSession
	result := dbConn.Where("username = ?", username).First(&foundUser)

	if result.Error == nil {
		client.UserSessions[chatID] = &models.UserSession{
			Id:       chatID,
			Step:     "done",
			Username: foundUser.Username,
			Name:     foundUser.Name,
			Email:    foundUser.Email,
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
			Id:       chatID,
			Step:     "name",
			Username: update.Message.Chat.Username,
		}
		client.SendMessageWithButtons(chatID, "Enter your name", []string{"Cancel"}, []string{"cancel"})
		return nil
	}

	session := client.UserSessions[chatID]

	switch step {
	case "start":
		client.SendMessage(chatID, "Enter your name")
	case "name":
		session.Name = text
		session.Step = "email"
		client.SendMessage(chatID, "Enter your email")
	case "email":
		session.Email = text
		if !isValidEmail(session.Email) {
			return ErrorHelper(err, "Email not valid")
		}
		session.Step = "done"
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
