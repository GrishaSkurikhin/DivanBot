package bothandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(log logger.BotLogger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var (
			handler            = "Start"
			username, inputMsg = requestinfo.Get(update)
		)

		messagesender.InfoWithKeyboard(ctx, b, update, log, handler, username, inputMsg,
			"Hello, пидарасы!!!", keyboards.MainMenu())
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
