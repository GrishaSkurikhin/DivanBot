package bothandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type AboutInfoGetter interface {
	GetAboutInfo() (string, error)
}

func About(log logger.BotLogger, aboutInfoGetter AboutInfoGetter) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler            = "About"
			username, inputMsg = requestinfo.Get(update)
		)

		info, err := aboutInfoGetter.GetAboutInfo()
		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get info", err)
			return
		}

		messagesender.Info(ctx, b, update, log, handler, username, inputMsg, info)
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
