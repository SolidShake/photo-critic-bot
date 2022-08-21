package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Factory struct {
	builder Builder
}

func NewFactory(builder Builder) Factory {
	return Factory{builder: builder}
}

func (f Factory) HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	switch message.Text {
	case sendInstaButtonText:
		msg = f.builder.SendInstagramButton(message.Chat.ID, message.Text)
	case reviewInstaButtonText:
		msg = f.builder.SendReviewInstagramButton(message.Chat.ID, message.Text)
	default:
		msg = f.builder.HandleUserMessage(message.Chat.ID, message.From.ID, message)
	}

	// @TODO add keyboard variants
	msg.ReplyMarkup = keyboard
	//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	return msg
}
