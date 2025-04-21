package tests_test

import (
	"telegram-bot/internal/handlers"
	"telegram-bot/internal/models"
	"telegram-bot/internal/repository"
	"testing"
)

const (
	loadPath  string = "/app/.env"
	fakeHost  string = "fake-host-telegram"
	fakeToken string = "fake-token-telegram"
)

func TestHandleCallback(t *testing.T) {
	update := models.Update{
		CallbackQuery: &models.CallbackQuery{
			Data: "help",
			Message: &models.Message{
				Chat: &models.Chat{
					Id: 111,
				},
			},
		},
	}
	client := repository.New(fakeHost, fakeToken)
	handlers.HandleCallback(update, client)
}

func TestHandleMessage(t *testing.T) {
	update := models.Update{
		Message: &models.Message{
			Text: "/start",
			Chat: &models.Chat{
				Id: 11,
			},
		},
	}

	client := repository.New(fakeHost, fakeToken)
	handlers.HandleMessage(update, client)
}
