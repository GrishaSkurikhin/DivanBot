package inlinehandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type FilmRegDeleter interface {
	DeleteRegOnFilm(userID uint64, filmID uint64) error
}

func CancelRegOnFilm(log logger.BotLogger, filmRegDeleter FilmRegDeleter, films []models.Film, userID int) slider.OnSelectFunc {
	return func(ctx context.Context, b *bot.Bot, message *botModels.Message, item int) {
		var (
			handler            = "CancelRegOnFilm"
			username, inputMsg = requestinfo.Get(&botModels.Update{Message: message})
		)

		err := filmRegDeleter.DeleteRegOnFilm(uint64(userID), films[item].ID)
		if err != nil {
			messagesender.Error(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to delete reg on film", err)
			return
		}

		messagesender.Info(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Запись успешно удалена")
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
