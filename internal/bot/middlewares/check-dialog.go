package middlewares

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func CheckDialog(log logger.BotLogger, d *dialoger.Dialoger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			var handler = "CheckDialog"
			var username, inputMsg string
			var chatID int64
			if update.CallbackQuery != nil {
				chatID  = update.CallbackQuery.Message.Chat.ID
				username = update.CallbackQuery.Message.From.Username
				inputMsg = update.CallbackQuery.Data
			} else {
				chatID  = update.Message.Chat.ID
				username = update.Message.From.Username
				inputMsg = update.Message.Text
			}
			
			dialogType, _, err := d.CheckDialog(chatID)
			if err != nil {
				log.BotERROR(handler, username, inputMsg, "Error on check dialog in middleware", err)
				messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
				return
			}

			if dialogType == dialoger.UnknownDialog {
				next(ctx, b, update)
			} else {
				if _, isExist := commands.Commands()[inputMsg]; isExist {
					messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Во время диалога команды не работают")
					if update.CallbackQuery != nil {
						b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
					}
					return
				}
				next(ctx, b, update)
			}
		}
	}
}

