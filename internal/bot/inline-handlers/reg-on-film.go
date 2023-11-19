package inlinehandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	customerrors "github.com/GrishaSkurikhin/DivanBot/internal/custom-errors"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type FilmRegistrator interface {
	RegOnFilm(userID uint64, filmID string) error
	IsUserReg(userID uint64) (bool, error)
	IsExistRegOnFilm(userID uint64, filmID string) (bool, error)
}

func RegOnFilm(log logger.BotLogger, filmRegistrator FilmRegistrator) slider.OnSelect {
	return func(ctx context.Context, b *bot.Bot, query *botModels.CallbackQuery, slideID string) {
		var (
			handler  = "RegOnFilm"
			username = query.Message.From.Username
			inputMsg = query.Message.Text
			chatID   = query.Message.Chat.ID
		)

		isReg, err := filmRegistrator.IsUserReg(uint64(chatID))
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to check user reg", err)
			return
		}

		if !isReg {
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Не зарегестрированы")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		isRegOnFilm, err := filmRegistrator.IsExistRegOnFilm(uint64(chatID), slideID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to check reg on film", err)
			return
		}

		if isRegOnFilm {
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Вы уже записаны на фильм")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		err = filmRegistrator.RegOnFilm(uint64(chatID), slideID)
		if err != nil {
			if _, ok := err.(customerrors.AlreadyRegistered); ok {
				messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Вы уже записаны на этот фильм")
				log.BotERROR(handler, username, inputMsg, "failed to reg on film", err)
			} else {
				messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
				log.BotERROR(handler, username, inputMsg, "failed to reg on film", err)
			}
			return
		}

		messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg,
			"Вы успешно записались")
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
