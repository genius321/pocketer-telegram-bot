package main

import (
	"github.com/boltdb/bolt"
	"github.com/genius321/pocketer-telegram-bot/internal/config"
	"github.com/genius321/pocketer-telegram-bot/internal/repository"
	"github.com/genius321/pocketer-telegram-bot/internal/repository/boltdb"
	"github.com/genius321/pocketer-telegram-bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/go-pocket-sdk"
)

const fsFileMode = 0600

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatal(err)
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

	db, err := initDB(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, cfg.RedirectURL, tokenRepository, cfg.Messages)

	if err := telegramBot.Start(); err != nil {
		logrus.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, fsFileMode, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
