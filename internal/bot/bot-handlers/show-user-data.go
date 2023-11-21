package bothandlers

import (
	"context"
	"fmt"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
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
			handler  = "ShowData"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
			userID   = uint64(update.Message.From.ID)
		)

		isReg, err := userDataGetter.IsUserReg(userID)
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

		name, surname, group, err := userDataGetter.GetUserData(userID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get user data", err)
			return
		}

		infoString := fmt.Sprintf("Ваши данные:\nИмя: %s\nФамилия: %s\nГруппа: %s", name, surname, group)
		kb := keyboards.ChangeData(b)

		messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username, inputMsg,
			infoString, kb)
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
