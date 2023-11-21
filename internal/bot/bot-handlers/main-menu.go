package bothandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

func MainMenu(log logger.BotLogger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler  = "MainMenu"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
		)

		messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username, inputMsg,
			"Главное меню", keyboards.MainMenu())
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
