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
	msg := tgbotapi.NewMessage(chatID, "Произошла неизвестная ошибка")

	switch err {
	case errInvalidURL:
		msg.Text = "Это невалидная ссылка!"
	case errUnableToSave:
		msg.Text = "Увы, не удалось сохранить ссылку. Попробуй ещё раз позже."
	case errAlreadyAuthorized:
		msg.Text = "Ты уже авторизирован. Присылай ссылку, а я её сохраню."
	case errUnknownCommand:
		msg.Text = "Я не знаю такой команды :("
	}
	b.bot.Send(msg)
}
