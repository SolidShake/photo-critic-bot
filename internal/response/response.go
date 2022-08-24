package response

import (
	"fmt"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	action "github.com/SolidShake/photo-critic-bot/internal/service/action"
	chat "github.com/SolidShake/photo-critic-bot/internal/service/chat"
	review "github.com/SolidShake/photo-critic-bot/internal/service/review"
)

const greetings = `Привет, %s! 👋

Это бот, в котором можно получить анонимную обратную связь по твоим фото-работам в Instagram.
`

type Builder struct {
	actionService *action.Service
	chatService   *chat.Service
	reviewService *review.Service
}

func NewBuilder(
	actionService *action.Service,
	chatService *chat.Service,
	reviewService *review.Service,
) Builder {
	return Builder{
		actionService: actionService,
		chatService:   chatService,
		reviewService: reviewService,
	}
}

func (b Builder) HandleUserMessage(chatID, fromID int64, message *tgbotapi.Message) tgbotapi.MessageConfig {
	action, _ := b.actionService.GetLastAction(chatID)

	switch action {
	case SendInstaButtonText:
		return b.SaveInstaResponse(chatID, message.Text)
	case ReviewInstaButtonText:
		return b.SaveInstaReviewResponse(chatID, message.Text) // insta reviewed id?
	case GetReviews:
		return b.GetReviewsResponse(chatID)
	default:
		return tgbotapi.NewMessage(chatID, b.DefaultMessage(chatID, message.Chat.FirstName))
	}
}

func (b Builder) SendInstagramButton(chatID int64, message string) tgbotapi.MessageConfig {
	err := b.actionService.SaveAction(chatID, message)
	if err != nil {
		return errorResponse(chatID, err)
	}

	return tgbotapi.NewMessage(chatID, "Отправь мне ссылку на свой Instagram \nМожно скопировать ее прямо из своего профиля")
}

func (b Builder) SendReviewInstagramButton(chatID int64, message string) tgbotapi.MessageConfig {
	chat, err := b.chatService.GetInstaLinkForReview(chatID)
	if err != nil {
		return errorResponse(chatID, err)
	}

	err = b.actionService.SaveAction(chatID, message)
	if err != nil {
		return errorResponse(chatID, err)
	}

	return tgbotapi.NewMessage(chatID, fmt.Sprintf("Instagram на оценку: %s \n\nВведите своё ревью профиля:", chat.Link))
}

func (b Builder) SaveInstaResponse(chatID int64, message string) tgbotapi.MessageConfig {
	message = strings.TrimSpace(message)

	if !isInstaLink(message) {
		// @TODO change description
		return tgbotapi.NewMessage(chatID, "Введите без пробелов")
	}

	message = addPrefixIfNeed(message)

	if err := b.chatService.SaveInstaLink(chatID, message); err != nil {
		return tgbotapi.NewMessage(chatID, "Не удалось сохранить ссылку, ошибка")
	}

	err := b.actionService.SaveAction(chatID, "default action")
	if err != nil {
		return errorResponse(chatID, err)
	}

	return tgbotapi.NewMessage(chatID, "Ссылка сохранена")
}

func (b Builder) SaveInstaReviewResponse(chatID int64, message string) tgbotapi.MessageConfig {
	chat, err := b.chatService.GetInstaLinkForReview(chatID)
	if err != nil {
		return errorResponse(chatID, err)
	}

	err = b.actionService.SaveAction(chatID, "default action")
	if err != nil {
		return errorResponse(chatID, err)
	}

	err = b.reviewService.SaveReview(chatID, chat.ChatID, message)
	if err != nil {
		return errorResponse(chatID, err)
	}

	return tgbotapi.NewMessage(chatID, "Спасибо за ревью! 💖 Ревью сохранено")
}

func (b Builder) GetReviewsResponse(chatID int64) tgbotapi.MessageConfig {
	err := b.actionService.SaveAction(chatID, "default action")
	if err != nil {
		return errorResponse(chatID, err)
	}

	reviews, err := b.reviewService.GetReviews(chatID)
	// @TODO handle empty
	if err != nil {
		return errorResponse(chatID, err)
	}

	var responseReviews string
	for _, review := range reviews {
		responseReviews += fmt.Sprintf(
			"<b>Ревью от %s</b>\n%s\n\n",
			review.GetFormatedTime(),
			review.Review,
		)
	}

	return tgbotapi.NewMessage(chatID, responseReviews)
}

func (b Builder) DefaultMessage(chatID int64, firstName string) string {
	message := fmt.Sprintf(greetings, firstName)

	link, err := b.chatService.GetInstaLink(chatID)
	if err == nil && link != "" {
		message += "\nТвой инстаграм: " + link
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

func errorResponse(chatID int64, err error) tgbotapi.MessageConfig {
	fmt.Println(err)
	return tgbotapi.NewMessage(chatID, "Упс, что-то пошло не так...")
}
