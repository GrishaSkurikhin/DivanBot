package bothandlers

import (
	"context"
	"fmt"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type UserDataGetter interface {
	GetUserData(userID uint64) (string, string, string, error) // name, surname, group, err
	IsUserReg(userID uint64) (bool, error)
}

func ShowData(log logger.BotLogger, userDataGetter UserDataGetter, d *dialoger.Dialoger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler            = "ShowData"
			username, inputMsg = requestinfo.Get(update)
			userID             = int(update.Message.From.ID)
		)

		isReg, err := userDataGetter.IsUserReg(uint64(userID))
		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to check user reg", err)
			return
		}

		if !isReg {
			messagesender.Info(ctx, b, update, log, handler, username, inputMsg, "Не зарегестрированы")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		name, surname, group, err := userDataGetter.GetUserData(uint64(userID))
		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get user data", err)
			return
		}

		infoString := fmt.Sprintf("Ваши данные:\nИмя: %s\nФамилия: %s\nГруппа: %s", name, surname, group)
		kb := keyboards.ChangeData(b)

		messagesender.InfoWithKeyboard(ctx, b, update, log, handler, username, inputMsg,
			infoString, kb)
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}