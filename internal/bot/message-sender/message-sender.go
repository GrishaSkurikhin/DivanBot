package messagesender

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Info(ctx context.Context, b *bot.Bot, chatID int64, log logger.BotLogger,
	handler, username, inputMsg, info string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   info,
	})

	if err != nil {
		log.BotERROR(handler, username, inputMsg, "Failed to send message", err)
	}
}

func InfoWithKeyboard(ctx context.Context, b *bot.Bot, chatID int64, log logger.BotLogger,
	handler, username, inputMsg, info string, keyboard models.ReplyMarkup) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        info,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		log.BotERROR(handler, username, inputMsg, "Failed to send message", err)
	}
}

func Error(ctx context.Context, b *bot.Bot, chatID int64, log logger.BotLogger,
	handler, username, inputMsg, errorMsg string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   errorMsg,
	})

	if err != nil {
		log.BotERROR(handler, username, inputMsg, "Failed to send error-message", err)
	}
}
