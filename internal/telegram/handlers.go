package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart = "start"
	startMessage = "Привет! Чтобы начать сохранять ссылки в своем Pocket аккаунте, для начала тебе необходимо дать мне на это доступ. Для этого переходи по ссылке:\n%s"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	logrus.Printf("[%s] %s", message.From.UserName, message.Text)
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return errUnknownCommand
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	logrus.Printf("[%s] %s", message.From.UserName, message.Text)

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}

	accessToken, err := b.getAccessTokenIfAuthorized(message.Chat.ID)
	if err != nil {
		return b.startAutorizationProcess(message.Chat.ID)
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		if err.Error() != "failed to parse response body: invalid semicolon separator in query" {
			return errUnableToSave
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена!")
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessTokenIfAuthorized(message.Chat.ID)
	if err != nil {
		return b.startAutorizationProcess(message.Chat.ID)
	}
	return errAlreadyAuthorized
}
