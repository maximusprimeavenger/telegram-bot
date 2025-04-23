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
	split := strings.Split(outputs[2], "\n")
	assert.Equal(t, "Unknown callback data: test", strings.TrimSuffix(outputs[1], "\n"), fmt.Sprintf("Not equal outputs:%s", outputs[1]))
	assert.Equal(t, "Empty callback query data from chat: 5", split[0], fmt.Sprintf("Not equal outputs:%s", outputs[2]))
}

func TestHandleMessage(t *testing.T) {
	updates := []models.Update{
		{
			Message: &models.Message{
				Text: "/start",
				Chat: &models.Chat{
					Id: 11,
				},
			},
		},
		{
			Message: &models.Message{
				Text: "test",
				Chat: &models.Chat{
					Id: 55,
				},
			},
		},
		{
			Message: &models.Message{
				Text: "",
				Chat: &models.Chat{
					Id: 5,
				},
			},
		},
	}
	for _, update := range updates {
		client := repository.New(fakeHost, fakeToken)
		handlers.HandleMessage(update, client)
	}
}
