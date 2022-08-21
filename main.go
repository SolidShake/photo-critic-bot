package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

const greetings = `Привет, %s! 👋

Это бот, в котором можно получить анонимную обратную связь по твоим фото-работам в Instagram.
`

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

var conn *pgx.Conn

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	conn, err = pgx.Connect(context.Background(), os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("cannot ping database: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("cannot start bot: %s", err)
	}

	bot.Debug = false // add to config

	updateConfig := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		var msg tgbotapi.MessageConfig

		switch update.Message.Text {
		case sendInstaButtonText:
			msg = sendInstagramButton(update.Message.Chat.ID, update.Message.From.ID, update.Message.Text)
		case reviewInstaButtonText:
			msg = sendReviewInstagramButton(update.Message.Chat.ID, update.Message.From.ID, update.Message.Text)
		default:
			msg = handleUserMessage(update.Message.Chat.ID, update.Message.From.ID, update.Message)
		}

		msg.ReplyMarkup = keyboard

		//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		// msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			log.Printf("send message error: %s", err)
		}
	}
}

func sendInstagramButton(chatID, fromID int64, message string) tgbotapi.MessageConfig {
	_ = saveAction(chatID, fromID, message)
	return tgbotapi.NewMessage(chatID, "Отправь мне ссылку на свой Instagram \nМожно скопировать ее прямо из своего профиля")
}

func sendReviewInstagramButton(chatID, fromID int64, message string) tgbotapi.MessageConfig {
	_ = saveAction(chatID, fromID, message)
	// @TODO get instagram link
	return tgbotapi.NewMessage(chatID, "Instagram на оценку: instagram.com/sofya.khvorostova/ \n\nВведите своё ревью профиля:")
}

func handleUserMessage(chatID, fromID int64, message *tgbotapi.Message) tgbotapi.MessageConfig {
	action, _ := getLastAction(chatID)
	switch action {
	case sendInstaButtonText:
		return saveInstaResponse(chatID, message.Text)
	case reviewInstaButtonText:
		return saveInstaReviewResponse(chatID, message.Text) // insta reviewed id?
	default:
		return tgbotapi.NewMessage(chatID, defaultMessage(chatID, message.Chat.FirstName))
	}
}

func saveInstaResponse(chatID int64, message string) tgbotapi.MessageConfig {
	message = strings.TrimSpace(message)

	if !isInstaLink(message) {
		// @TODO change description
		return tgbotapi.NewMessage(chatID, "Введите без пробелов")
	}

	message = addPrefixIfNeed(message)

	if err := saveInsta(chatID, message); err != nil {
		return tgbotapi.NewMessage(chatID, "Не удалось сохранить ссылку, ошибка")
	}

	// @TODO change desc
	_ = saveAction(chatID, chatID, "default action")
	return tgbotapi.NewMessage(chatID, "Ссылка сохранена")
}

func saveInstaReviewResponse(chatID int64, message string) tgbotapi.MessageConfig {
	// @TODO change desc
	_ = saveAction(chatID, chatID, "default action")
	return tgbotapi.NewMessage(chatID, "Спасибо за ревью! 💖 Ревью сохранено")
}

func defaultMessage(chatID int64, firstName string) string {
	message := fmt.Sprintf(greetings, firstName)

	link, err := getInsta(chatID)
	if err == nil && link != "" {
		message += "Твой инстаграм: " + link
	}

	return message
}

func saveAction(chatID, userID int64, action string) error {
	row := conn.QueryRow(context.Background(), "INSERT INTO actions (chat_id, user_id, action) VALUES ($1, $2, $3) RETURNING id", chatID, userID, action)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("insert db error: %s", err)
	}

	return nil
}

func saveInsta(chatID int64, link string) error {
	row := conn.QueryRow(context.Background(), "INSERT INTO links (chat_id, link) VALUES ($1, $2) RETURNING id", chatID, link)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("insert db error: %s", err)
	}

	return nil
}

func getInsta(chatID int64) (string, error) {
	var link string
	err := conn.QueryRow(context.Background(), "SELECT link FROM links WHERE chat_id = $1 ORDER BY id DESC", chatID).Scan(&link)
	if err != nil {
		return "", fmt.Errorf("insert db error: %s", err)
	}

	return link, nil
}

func getLastAction(chatID int64) (string, error) {
	var action string
	err := conn.QueryRow(context.Background(), "SELECT action FROM actions WHERE chat_id = $1 ORDER BY id DESC", chatID).Scan(&action)
	if err != nil {
		return "", fmt.Errorf("insert db error: %s", err)
	}

	return action, nil
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
