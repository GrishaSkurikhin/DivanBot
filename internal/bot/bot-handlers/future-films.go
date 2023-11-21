package bothandlers

import (
	"context"

	inlinehandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/inline-handlers"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/sliders"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type FutureFilmsGetter interface {
	GetFutureFims() ([]models.Film, error)
	inlinehandlers.FilmRegistrator
}

func FutureFilms(log logger.BotLogger, futureFilmsGetter FutureFilmsGetter) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler  = "FutureFilms"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
		)

		films, err := futureFilmsGetter.GetFutureFims()
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get future films", err)
			return
		}
		if len(films) == 0 {
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Пока нет запланированных фильмов")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		sl := sliders.FutureFilms(films)
		_, err = sl.Show(ctx, b, chatID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
