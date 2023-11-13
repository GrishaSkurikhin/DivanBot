package bothandlers

import (
	"context"
	"fmt"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

func Default(log logger.BotLogger, d *dialoger.Dialoger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		if update.Message == nil {
			// fmt.Println(update.CallbackQuery.Message)
			fmt.Println(update.CallbackQuery.Data)
			return
		}

		var (
			handler            = "Default"
			username, inputMsg = requestinfo.Get(update)
			chatID             = int(update.Message.Chat.ID)
		)

		dialogType, state, err := d.CheckDialog(chatID)

		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to check dialog", err)
			return
		}

		if dialogType != dialoger.UnknownDialog {
			err := d.ServeMessage(ctx, b, update, dialogType, state)
			if err != nil {
				messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
				log.BotERROR(handler, username, inputMsg, "Failed to serve message", err)
				return
			}

			if inputMsg == commands.Cancel {
				messagesender.InfoWithKeyboard(ctx, b, update, log, handler, username,
					inputMsg, "Операция отменена", keyboards.MainMenu())
			}
			return
		}

		messagesender.Info(ctx, b, update, log, handler, username, inputMsg, "Неизвестная команда")
	}
}
