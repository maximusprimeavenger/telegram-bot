package tests_test

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"telegram-bot/internal/handlers"
	"telegram-bot/internal/models"
	"telegram-bot/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	loadPath  string = "/app/.env"
	fakeHost  string = "fake-host-telegram"
	fakeToken string = "fake-token-telegram"
)

/*
queryChatId := update.CallbackQuery.Message.Chat.Id

	data := update.CallbackQuery.Data
*/
func TestHandleCallback(t *testing.T) {
	updates := []models.Update{
		{
			CallbackQuery: &models.CallbackQuery{
				Data: "help",
				Message: &models.Message{
					Chat: &models.Chat{
						Id: 111,
					},
				},
			},
		},
		{
			CallbackQuery: &models.CallbackQuery{
				Data: "test",
				Message: &models.Message{
					Chat: &models.Chat{
						Id: 55,
					},
				},
			},
		},
		{
			CallbackQuery: &models.CallbackQuery{
				Data: "",
				Message: &models.Message{
					Chat: &models.Chat{
						Id: 5,
					},
				},
			},
		},
	}
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	outputs := make([]string, 3)
	client := repository.New(fakeHost, fakeToken)
	for i, update := range updates {
		buf.Reset()
		handlers.HandleCallback(update, client)
		outputs[i] = buf.String()
	}

	assert.Equal(t, "Unknown callback data: test", strings.TrimSuffix(outputs[1], "\n"), fmt.Sprintf("Not equal outputs:%s", outputs[1]))
	assert.Equal(t, "Empty callback query data from chat: 5\nUnknown callback data: \n", outputs[2], fmt.Sprintf("Not equal outputs:%s", outputs[2]))
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
