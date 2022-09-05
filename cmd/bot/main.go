package main

import (
	"os"

	"github.com/genius321/pocketer-telegram-bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Println(err.Error())
		return
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		logrus.Fatal(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot)
	if err := telegramBot.Start(); err != nil {
		logrus.Fatal(err)
	}
}
