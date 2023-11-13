package bothandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

const (
	helpInfo = "пиздатая помощь"
)

func Help(log logger.BotLogger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler            = "Help"
			username, inputMsg = requestinfo.Get(update)
		)

		messagesender.Info(ctx, b, update, log, handler, username, inputMsg, helpInfo)
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
