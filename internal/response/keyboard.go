package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	SendInstaButtonText   = "Отправить свой Instagram на оценку"
	ReviewInstaButtonText = "Оценить чужой Instagram"
)

var Keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(SendInstaButtonText),
		tgbotapi.NewKeyboardButton(ReviewInstaButtonText),
	),
)
