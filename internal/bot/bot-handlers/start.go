package bothandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	startInfo = `
Привет!

Это телеграм-бот киноклуба Диван. Здесь ты сможешь зарегистрироваться на будущие киновечера или посмотреть, какие фильмы диван уже показал.

Информация по управлению ботом доступна по команде /help. 
`
)

func Start(log logger.BotLogger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var (
			handler  = "Start"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID = update.Message.Chat.ID
		)

		messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username, inputMsg, startInfo, keyboards.MainMenu())
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
