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

type FilmsRegsGetter interface {
	GetFilmsRegs(userID uint64) ([]models.Film, error)
	inlinehandlers.FilmRegDeleter
	IsUserReg(userID uint64) (bool, error)
}

func ShowRegs(log logger.BotLogger, filmsRegsGetter FilmsRegsGetter) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler  = "ShowRegs"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			userID   = uint64(update.Message.Chat.ID)
			chatID   = update.Message.Chat.ID
		)

		isReg, err := filmsRegsGetter.IsUserReg(userID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to check user reg", err)
			return
		}

		if !isReg {
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Вы не зарегистрированы. Для регистрации введите /reg")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		films, err := filmsRegsGetter.GetFilmsRegs(userID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get previous films", err)
			return
		}

		if len(films) == 0 {
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "У вас нет записей на фильмы")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		sl := sliders.RegFilms(films)
		_, err = sl.Show(ctx, b, chatID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}