package messagesender

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Info(ctx context.Context, b *bot.Bot, update *models.Update, log logger.BotLogger,
	handler, username, inputMsg, info string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   info,
	})

	if err != nil {
		log.BotERROR(handler, username, inputMsg, "Failed to send message", err)
	}
}

func InfoWithKeyboard(ctx context.Context, b *bot.Bot, update *models.Update, log logger.BotLogger,
	handler, username, inputMsg, info string, keyboard models.ReplyMarkup) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        info,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		log.BotERROR(handler, username, inputMsg, "Failed to send message", err)
	}
}

func Error(ctx context.Context, b *bot.Bot, update *models.Update, log logger.BotLogger,
	handler, username, inputMsg, errorMsg string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   errorMsg,
	})

	if err != nil {
		log.BotERROR(handler, username, inputMsg, "Failed to send error-message", err)
	}
}
