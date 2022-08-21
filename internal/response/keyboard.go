package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	sendInstaButtonText   = "Отправить свой Instagram на оценку"
	reviewInstaButtonText = "Оценить чужой Instagram"
)

var keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(sendInstaButtonText),
		tgbotapi.NewKeyboardButton(reviewInstaButtonText),
	),
)
