package inlinehandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

func ShowLocation(log logger.BotLogger, films []models.Film) slider.OnSelectFunc {
	return func(ctx context.Context, b *bot.Bot, message *botModels.Message, item int) {
		var (
			handler            = "ShowLocation"
			username, inputMsg = requestinfo.Get(&botModels.Update{Message: message})
			chatID             = int(message.Chat.ID)
		)
		loc := films[item].Location

		_, err := b.SendVenue(ctx, &bot.SendVenueParams{
			ChatID:    chatID,
			Latitude:  loc.Lat,
			Longitude: loc.Long,
			Title:     loc.Title,
			Address:   loc.Description,
		})
		if err != nil {
			messagesender.Error(ctx, b, &botModels.Update{Message: message}, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
