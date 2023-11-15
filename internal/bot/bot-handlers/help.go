package bothandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

const (
	helpInfo = "помощь"
)

func Help(log logger.BotLogger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler  = "Help"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
		)

		messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, helpInfo)
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
