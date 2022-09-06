package telegram

import (
	"context"
	"fmt"

	"github.com/genius321/pocketer-telegram-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) getAccessTokenFromDB(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) getRequestTokenFromDB(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.RequestTokens)
}

func (b *Bot) startAutorizationProcess(chatID int64) error {
	authLink, err := b.generateAutorizationLink(chatID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf(startMessage, authLink))

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) generateAutorizationLink(chatID int64) (string, error) {
	redirectURL := b.generateRedirectURL(chatID)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectURL(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
