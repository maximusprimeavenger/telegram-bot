package models

type Update struct {
	Id            int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Message struct {
	Text        string      `json:"text"`
	Chat        *Chat       `json:"chat"`
	ReplyMarkup interface{} `json:"reply_markup,omitempty"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type CallbackQuery struct {
	Id      string        `json:"id"`
	From    *TelegramUser `json:"from"`
	Message *Message      `json:"message,omitempty"`
	Data    string        `json:"data"`
}

type TelegramUser struct {
	Id       int    `json:"id"`
	Username string `json:"username,omitempty"`
}

// For taking response from buttons
type InlineKeyboard struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}
