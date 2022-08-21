package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"

	actionRepository "github.com/SolidShake/photo-critic-bot/internal/repository/action"
	linkRepository "github.com/SolidShake/photo-critic-bot/internal/repository/link"
	"github.com/SolidShake/photo-critic-bot/internal/response"
	actionService "github.com/SolidShake/photo-critic-bot/internal/service/action"
	linkService "github.com/SolidShake/photo-critic-bot/internal/service/link"
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

	factory := response.NewFactory(
		response.NewBuilder(
			actionService.NewActionService(
				actionRepository.NewRepository(conn),
			),
			linkService.NewLinkService(
				linkRepository.NewRepository(conn),
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
