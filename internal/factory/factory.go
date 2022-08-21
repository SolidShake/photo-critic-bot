package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	response "github.com/SolidShake/photo-critic-bot/internal/response"
)

type Factory struct {
	builder response.Builder
}

func NewFactory(builder response.Builder) Factory {
	return Factory{builder: builder}
}

func (f Factory) HandleMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	switch message.Text {
	case response.SendInstaButtonText:
		msg = f.builder.SendInstagramButton(message.Chat.ID, message.Text)
	case response.ReviewInstaButtonText:
		msg = f.builder.SendReviewInstagramButton(message.Chat.ID, message.Text)
	default:
		msg = f.builder.HandleUserMessage(message.Chat.ID, message.From.ID, message)
	}

	// @TODO add keyboard variants
	msg.ReplyMarkup = response.Keyboard
	//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	return msg
}
