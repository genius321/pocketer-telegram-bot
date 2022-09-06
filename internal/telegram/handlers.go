package telegram

import (
	"context"
	"net/url"

	"github.com/genius321/pocketer-telegram-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart           = "start"
	startMessage           = "Привет! Чтобы начать сохранять ссылки в своем Pocket аккаунте, для начала тебе необходимо дать мне на это доступ. Для этого переходи по ссылке:\n%s"
	replyAlreadyAuthorized = "Ты уже авторизирован. Присылай ссылку, а я её сохраню."
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	// logrus.Printf("[%s] %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Ссылка успешно сохранена!")

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "Это невалидная ссылка!"
		_, err := b.bot.Send(msg)
		return err
	}

	accessToken, err := b.getAccessTokenFromDB(message.Chat.ID)
	if err != nil {
		requestToken, err := b.getRequestTokenFromDB(message.Chat.ID)
		if err != nil {
			msg.Text = "Ты не авторизирован! Используй команду /start"
			_, err := b.bot.Send(msg)
			return err
		}

		authResponse, err := b.pocketClient.Authorize(context.Background(), requestToken)
		if err != nil {
			msg.Text = "Ты не авторизирован! Используй команду /start"
			_, err := b.bot.Send(msg)
			return err
		}

		if err := b.tokenRepository.Save(message.Chat.ID, authResponse.AccessToken, repository.AccessTokens); err != nil {
			msg.Text = "Ты не авторизирован! Используй команду /start"
			_, err := b.bot.Send(msg)
			return err
		}
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
	_, err := b.getAccessTokenFromDB(message.Chat.ID)
	if err != nil {
		requestToken, err := b.getRequestTokenFromDB(message.Chat.ID)
		if err != nil {
			return b.startAutorizationProcess(message.Chat.ID)
		}

		// get access token on pocket
		authResponse, err := b.pocketClient.Authorize(context.Background(), requestToken)
		if err != nil {
			return b.startAutorizationProcess(message.Chat.ID)
		}

		if err := b.tokenRepository.Save(message.Chat.ID, authResponse.AccessToken, repository.AccessTokens); err != nil {
			return b.startAutorizationProcess(message.Chat.ID)
		}
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
