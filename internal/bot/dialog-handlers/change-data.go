package dialoghandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type UserDataChanger interface {
	ChangeUserData(dataType string, newValue string, userID uint64) error
}

func ChangeData(log logger.BotLogger, userDataChanger UserDataChanger) dialoger.DialogHandler {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update, state int, info map[string]string) (newInfo map[string]string, isErr bool) {
		var (
			handler            = "ChangeData"
			username, inputMsg = requestinfo.Get(update)
		)
		newInfo = make(map[string]string)

		switch state {
		case 1:
			messagesender.Info(ctx, b, update, log, handler, username, inputMsg, "Введите новое значение")
		case 2:
			dataType := info["dataType"]
			data := inputMsg

			err := userDataChanger.ChangeUserData(dataType, data, uint64(update.Message.From.ID))
			if err != nil {
				messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
				log.BotERROR(handler, username, inputMsg, "Failed to change info in db", err)
				isErr = true
				return
			}

			messagesender.InfoWithKeyboard(ctx, b, update, log, handler, username, inputMsg,
				"Информация успешно изменена", keyboards.MainMenu())
		}

		log.BotINFO(handler, username, inputMsg, "successfully")
		return
	}
}
