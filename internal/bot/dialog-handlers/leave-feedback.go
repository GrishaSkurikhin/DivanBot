package dialoghandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type FeedbackSender interface {
	SendFeedback(userID uint64, comment string) error
}

func LeaveFeedback(log logger.BotLogger, feedbackSender FeedbackSender) dialoger.DialogHandler {
	return func(ctx context.Context, b *bot.Bot, msg *botModels.Message, state int, info map[string]string) (newInfo map[string]string, isErr bool) {
		var (
			handler  = "LeaveFeedback"
			username = msg.From.Username
			inputMsg = msg.Text
			userID   = uint64(msg.From.ID)
			chatID = msg.Chat.ID
		)
		newInfo = make(map[string]string)

		switch state {
		case 1:
			messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username,
				inputMsg, "Введите ваши пожелания и (или) отзыв", keyboards.DialogMenu())
		case 2:
			feedback := inputMsg
			err := feedbackSender.SendFeedback(userID, feedback)
			if err != nil {
				messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
				log.BotERROR(handler, username, inputMsg, "Failed to add feedback to db", err)
				isErr = true
				return
			}

			messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username, inputMsg,
				"Спасибо за отзыв!", keyboards.MainMenu())
		}

		log.BotINFO(handler, username, inputMsg, "successfully")
		return
	}
}
