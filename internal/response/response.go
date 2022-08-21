package response

import (
	"fmt"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	action "github.com/SolidShake/photo-critic-bot/internal/service/action"
	link "github.com/SolidShake/photo-critic-bot/internal/service/link"
)

const greetings = `Привет, %s! 👋

Это бот, в котором можно получить анонимную обратную связь по твоим фото-работам в Instagram.
`

type Builder struct {
	actionService *action.ActionService
	linkService   *link.LinkService
}

func NewBuilder(actionService *action.ActionService, linkService *link.LinkService) Builder {
	return Builder{
		actionService: actionService,
		linkService:   linkService,
	}
}

func (b Builder) HandleUserMessage(chatID, fromID int64, message *tgbotapi.Message) tgbotapi.MessageConfig {
	action, _ := b.actionService.GetLastAction(chatID)

	switch action {
	case sendInstaButtonText:
		return b.SaveInstaResponse(chatID, message.Text)
	case reviewInstaButtonText:
		return b.SaveInstaReviewResponse(chatID, message.Text) // insta reviewed id?
	default:
		return tgbotapi.NewMessage(chatID, b.DefaultMessage(chatID, message.Chat.FirstName))
	}
}

func (b Builder) SendInstagramButton(chatID int64, message string) tgbotapi.MessageConfig {
	_ = b.actionService.SaveAction(chatID, message)
	return tgbotapi.NewMessage(chatID, "Отправь мне ссылку на свой Instagram \nМожно скопировать ее прямо из своего профиля")
}

func (b Builder) SendReviewInstagramButton(chatID int64, message string) tgbotapi.MessageConfig {
	_ = b.actionService.SaveAction(chatID, message)
	// @TODO get instagram link
	return tgbotapi.NewMessage(chatID, "Instagram на оценку: instagram.com/sofya.khvorostova/ \n\nВведите своё ревью профиля:")
}

func (b Builder) SaveInstaResponse(chatID int64, message string) tgbotapi.MessageConfig {
	message = strings.TrimSpace(message)

	if !isInstaLink(message) {
		// @TODO change description
		return tgbotapi.NewMessage(chatID, "Введите без пробелов")
	}

	message = addPrefixIfNeed(message)

	if err := b.linkService.SaveInstaLink(chatID, message); err != nil {
		return tgbotapi.NewMessage(chatID, "Не удалось сохранить ссылку, ошибка")
	}

	// @TODO change desc
	_ = b.actionService.SaveAction(chatID, "default action")
	return tgbotapi.NewMessage(chatID, "Ссылка сохранена")
}

func (b Builder) SaveInstaReviewResponse(chatID int64, message string) tgbotapi.MessageConfig {
	// @TODO change desc
	_ = b.actionService.SaveAction(chatID, "default action")
	return tgbotapi.NewMessage(chatID, "Спасибо за ревью! 💖 Ревью сохранено")
}

func (b Builder) DefaultMessage(chatID int64, firstName string) string {
	message := fmt.Sprintf(greetings, firstName)

	link, err := b.linkService.GetInstaLink(chatID)
	if err == nil && link != "" {
		message += "Твой инстаграм: " + link
	}

	return message
}

func isInstaLink(text string) bool {
	// @TODO add insta verification
	m, _ := regexp.MatchString("^\\S+$", text)
	return m
}

func addPrefixIfNeed(text string) string {
	if strings.Contains(text, "instagram.com") {
		return text
	}

	return "https://instagram.com/" + text
}
