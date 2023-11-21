package bothandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type IsUserRegChecker interface {
	IsUserReg(userID uint64) (bool, error)
}

func LeaveFeedbackStart(log logger.BotLogger, d *dialoger.Dialoger, isUserRegChecker IsUserRegChecker) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var (
			handler  = "LeaveFeedback"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
			userID   = uint64(update.Message.From.ID)
		)

		isReg, err := isUserRegChecker.IsUserReg(userID)
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

		err = d.StartDialog(ctx, b, update.Message, dialoger.LeaveFeedbackDialog, chatID, nil)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to start dialog", err)
		}
	}
}
