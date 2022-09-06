package main

import (
	"log"

	"github.com/genius321/pocketer-telegram-bot/internal/config"
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

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		logrus.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.ConsumerKey)
	if err != nil {
		logrus.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, pocketClient, cfg.RedirectURL)
	if err := telegramBot.Start(); err != nil {
		logrus.Fatal(err)
	}
}
