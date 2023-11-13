package inlinehandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	inlinekeyboard "github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/inline-keyboard"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

func ChangeDataStart(log logger.BotLogger, d *dialoger.Dialoger) inlinekeyboard.OnSelect {
	return func(ctx context.Context, b *bot.Bot, mes *botModels.Message, data []byte) {
		var (
			handler            = "ChangeData"
			username, inputMsg = requestinfo.Get(&botModels.Update{Message: mes})
			chatID             = int(mes.Chat.ID)
		)

		startInfo := map[string]string{"dataType": string(data)}

		err := d.StartDialog(ctx, b, &botModels.Update{Message: mes}, dialoger.ChangeDataDialog, chatID, startInfo)
		if err != nil {
			messagesender.Error(ctx, b, &botModels.Update{Message: mes}, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to start dialog", err)
		}
	}
}
