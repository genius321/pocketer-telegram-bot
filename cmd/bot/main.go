package main

import (
	"os"

	"github.com/genius321/pocketer-telegram-bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		logrus.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(os.Getenv("CONSUMER_KEY"))
	if err != nil {
		logrus.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, pocketClient)
	if err := telegramBot.Start(); err != nil {
		logrus.Fatal(err)
	}
}
