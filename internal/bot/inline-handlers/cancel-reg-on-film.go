package inlinehandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type FilmRegDeleter interface {
	DeleteRegOnFilm(userID uint64, filmID uint64) error
}

func CancelRegOnFilm(log logger.BotLogger, filmRegDeleter FilmRegDeleter, films []models.Film, userID uint64) slider.OnSelect {
	return func(ctx context.Context, b *bot.Bot, query *botModels.CallbackQuery, item int) {
		var (
			handler  = "CancelRegOnFilm"
			username = query.Message.From.Username
			inputMsg = query.Message.Text
			chatID = query.Message.Chat.ID
		)

		err := filmRegDeleter.DeleteRegOnFilm(userID, films[item].ID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to delete reg on film", err)
			return
		}

		messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Запись успешно удалена")
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
