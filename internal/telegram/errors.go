package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errInvalidURL        = errors.New("url is invalid")
	errUnableToSave      = errors.New("unable to save")
	errAlreadyAuthorized = errors.New("already authorized")
	errUnknownCommand    = errors.New("unknown command")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Default)

	switch err {
	case errInvalidURL:
		msg.Text = b.messages.InvalidURL
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
	case errAlreadyAuthorized:
		msg.Text = b.messages.AlreadyAuthorized
	case errUnknownCommand:
		msg.Text = b.messages.UnknownCommand
	}
	b.bot.Send(msg)
}
