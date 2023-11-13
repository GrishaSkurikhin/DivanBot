package inlinehandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	customerrors "github.com/GrishaSkurikhin/DivanBot/internal/custom-errors"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type FilmRegistrator interface {
	RegOnFilm(userID uint64, filmID uint64) error
	IsUserReg(userID uint64) (bool, error)
	IsExistRegOnFilm(userID uint64, filmID uint64) (bool, error)
}

func RegOnFilm(log logger.BotLogger, filmRegistrator FilmRegistrator, films []models.Film, userID int) slider.OnSelectFunc {
	return func(ctx context.Context, b *bot.Bot, message *botModels.Message, item int) {
		var (
			handler            = "RegOnFilm"
			username, inputMsg = requestinfo.Get(&botModels.Update{Message: message})
		)

		isReg, err := filmRegistrator.IsUserReg(uint64(userID))
		if err != nil {
			messagesender.Error(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to check user reg", err)
			return
		}

		if !isReg {
			messagesender.Info(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Не зарегестрированы")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		isRegOnFilm, err := filmRegistrator.IsExistRegOnFilm(uint64(userID), films[item].ID)
		if err != nil {
			messagesender.Error(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to check reg on film", err)
			return
		}

		if isRegOnFilm {
			messagesender.Info(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Вы уже записаны на фильм")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		err = filmRegistrator.RegOnFilm(uint64(userID), films[item].ID)
		if err != nil {
			if _, ok := err.(customerrors.AlreadyRegistered); ok {
				messagesender.Error(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Вы уже записаны на этот фильм")
				log.BotERROR(handler, username, inputMsg, "failed to reg on film", err)
			} else {
				messagesender.Error(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Ошибка")
				log.BotERROR(handler, username, inputMsg, "failed to reg on film", err)
			}
			return
		}

		messagesender.Info(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg,
			"Вы успешно записались")
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
