package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"telegram-bot/internal/helpers"
	"telegram-bot/internal/models"
)

type Client struct {
	host         string
	basePath     string
	client       http.Client
	UserSessions map[int]*models.UserSession
}

const (
	getUpdates  = "getUpdates"
	sendMessage = "sendMessage"
)

func New(host, token string) *Client {
	return &Client{
		host:         host,
		basePath:     newBasePath(token),
		client:       http.Client{},
		UserSessions: make(map[int]*models.UserSession),
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]models.Update, error) {
	params := map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, helpers.ErrorHelper(err, "error marshaling parameters")
	}
	data, err := c.sendRequest(getUpdates, paramsJSON)
	if err != nil {
		return nil, helpers.ErrorHelper(err, "cannot read the data from body")
	}

	var res models.UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, helpers.ErrorHelper(err, "error unparsing json")
	}

	return res.Result, nil
}

func (c *Client) sendRequest(method string, data []byte) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, helpers.ErrorHelper(err, fmt.Sprint("can not create a request with url:", u))
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, helpers.ErrorHelper(err, "cannot send a request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, helpers.ErrorHelper(err, "cannot read the response")
	}

	fmt.Println("Telegram API response:", string(body))
	return body, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	data := map[string]interface{}{
		"chat_id": chatId,
		"text":    text,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return helpers.ErrorHelper(err, "error marshaling markup")
	}
	fmt.Println("Reply markup JSON:", string(dataJSON))
	_, err = c.sendRequest(sendMessage, dataJSON)
	if err != nil {
		return helpers.ErrorHelper(err, "error while sending a request")
	}
	return nil
}

func (c *Client) SendMessageWithButtons(chatId int, text string, params, callbackdata []string) error {
	inlineKeyboard := [][]models.InlineKeyboardButton{}
	for i := 0; i < len(params); i++ {
		inlineKeyboard = append(inlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         params[i],
				CallbackData: callbackdata[i],
			},
		})
	}

	body := map[string]interface{}{
		"chat_id": chatId,
		"text":    text,
		"reply_markup": models.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return helpers.ErrorHelper(err, "error marshaling body")
	}
	resp, err := c.sendRequest(sendMessage, bodyJSON)
	if err != nil {
		return helpers.ErrorHelper(err, "error sending a request")
	}
	log.Println("Response:", string(resp))
	return nil
}

func (c *Client) SendWithKeyboard(chatId int, text string, textButton []string) error {
	keyboardButton := [][]models.KeyboardButton{}
	for i := 0; i < len(textButton); i++ {
		keyboardButton = append(keyboardButton, []models.KeyboardButton{
			{Text: textButton[i]},
		})
	}
	keyboard := models.ReplyKeyboardMarkup{
		Keyboard:        keyboardButton,
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
	data := map[string]interface{}{
		"reply_markup": keyboard,
		"chat_id":      chatId,
		"text":         text,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return helpers.ErrorHelper(err, "error marshaling body")
	}
	resp, err := c.sendRequest(sendMessage, dataJSON)
	if err != nil {
		return helpers.ErrorHelper(err, "error sending a request")
	}
	log.Println("Response:", string(resp))
	return nil
}
