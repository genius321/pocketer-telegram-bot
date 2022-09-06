package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart           = "start"
	startMessage           = "Привет! Чтобы начать сохранять ссылки в своем Pocket аккаунте, для начала тебе необходимо дать мне на это доступ. Для этого переходи по ссылке:\n%s"
	replyAlreadyAuthorized = "Ты уже авторизирован. Присылай ссылку, а я её сохраню."
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	logrus.Printf("[%s] %s", message.From.UserName, message.Text)
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	logrus.Printf("[%s] %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена!")

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "Это невалидная ссылка!"
		_, err := b.bot.Send(msg)
		return err
	}

	accessToken, err := b.getAccessTokenIfAuthorized(message.Chat.ID)
	if err != nil {
		return b.startAutorizationProcess(message.Chat.ID)
	}

	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		msg.Text = "Не удалось сохранить ссылку. Попробуй ещё раз позже."
		_, err := b.bot.Send(msg)
		return err
	}

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessTokenIfAuthorized(message.Chat.ID)
	if err != nil {
		return b.startAutorizationProcess(message.Chat.ID)
	}

	// user authorized already
	msg := tgbotapi.NewMessage(message.Chat.ID, replyAlreadyAuthorized)

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды :(")

	_, err := b.bot.Send(msg)
	return err
}
