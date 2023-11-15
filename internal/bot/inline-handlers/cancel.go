package inlinehandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	inlinekeyboard "github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/inline-keyboard"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

func Cancel(log logger.BotLogger) inlinekeyboard.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, query *botModels.CallbackQuery) {
		var (
			handler  = "Cancel"
			username = query.Message.From.Username
			inputMsg = query.Data
		)

		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
