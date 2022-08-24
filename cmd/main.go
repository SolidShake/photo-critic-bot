package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"

	factory "github.com/SolidShake/photo-critic-bot/internal/factory"
	actionRepository "github.com/SolidShake/photo-critic-bot/internal/repository/action"
	chatRepository "github.com/SolidShake/photo-critic-bot/internal/repository/chat"
	reviewRepository "github.com/SolidShake/photo-critic-bot/internal/repository/review"
	"github.com/SolidShake/photo-critic-bot/internal/response"
	actionService "github.com/SolidShake/photo-critic-bot/internal/service/action"
	chatService "github.com/SolidShake/photo-critic-bot/internal/service/chat"
	reviewService "github.com/SolidShake/photo-critic-bot/internal/service/review"
)

// @TODO добавить проверку на добавление инициализируего профиля
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_HOST"))
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

	factory := factory.NewFactory(
		response.NewBuilder(
			actionService.NewService(
				actionRepository.NewRepository(conn),
			),
			chatService.NewService(
				chatRepository.NewRepository(conn),
			),
			reviewService.NewService(
				reviewRepository.NewRepository(conn),
			),
		),
	)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := factory.HandleMessage(update.Message)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("send message error: %s", err)
		}
	}
}
