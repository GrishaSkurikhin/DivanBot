package bothandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/sliders"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type PrevFilmsGetter interface {
	GetPrevFims() ([]models.Film, error)
}

func PrevFilms(log logger.BotLogger, prevFilmsGetter PrevFilmsGetter) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler  = "PrevFilms"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
		)

		films, err := prevFilmsGetter.GetPrevFims()
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get previous films", err)
			return
		}
		if len(films) == 0 {
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Пока не было показано ни одного фильма")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		sl := sliders.PrevFilms(films)
		_, err = sl.Show(ctx, b, chatID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
