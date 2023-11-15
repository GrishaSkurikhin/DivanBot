package inlinehandlers

import (
	"context"
	"strings"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	inlinekeyboard "github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/inline-keyboard"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

func ChangeDataStart(log logger.BotLogger, d *dialoger.Dialoger) inlinekeyboard.OnSelect {
	return func(ctx context.Context, bot *bot.Bot, query *botModels.CallbackQuery) {
		var (
			handler  = "ChangeData"
			username = query.Message.From.Username
			chatID   = query.Message.Chat.ID
			data     = string(query.Data)
		)

		data = strings.TrimPrefix(data, "data")
		startInfo := map[string]string{"dataType": data}

		err := d.StartDialog(ctx, bot, query.Message, dialoger.ChangeDataDialog, chatID, startInfo)
		if err != nil {
			messagesender.Error(ctx, bot, chatID, log, handler, username, data, "Ошибка")
			log.BotERROR(handler, username, data, "failed to start dialog", err)
		}
	}
}
