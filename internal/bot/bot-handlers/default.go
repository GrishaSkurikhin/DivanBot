package bothandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

func Default(log logger.BotLogger, d *dialoger.Dialoger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		if update.CallbackQuery != nil {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
			return
		}

		var (
			handler  = "Default"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID = update.Message.Chat.ID
		)

		dialogType, state, err := d.CheckDialog(chatID)

		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to check dialog", err)
			return
		}

		if dialogType != dialoger.UnknownDialog {
			err := d.ServeMessage(ctx, b, update.Message, dialogType, state)
			if err != nil {
				messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
				log.BotERROR(handler, username, inputMsg, "Failed to serve message", err)
				return
			}

			if inputMsg == commands.Cancel {
				messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username,
					inputMsg, "Операция отменена", keyboards.MainMenu())
			}
			return
		}

		messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Неизвестная команда")
	}
}
