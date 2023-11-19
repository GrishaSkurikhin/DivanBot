package inlinehandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type LocationGetter interface {
	GetLocation(filmID string) (models.Location, error)
}

func ShowLocation(log logger.BotLogger, locationGetter LocationGetter) slider.OnSelect {
	return func(ctx context.Context, b *bot.Bot, query *botModels.CallbackQuery, slideID string) {
		var (
			handler  = "ShowLocation"
			username = query.Message.From.Username
			inputMsg = query.Message.Text
			chatID   = query.Message.Chat.ID
		)
		loc, err := locationGetter.GetLocation(slideID)
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get location", err)
			return
		}

		_, err = b.SendVenue(ctx, &bot.SendVenueParams{
			ChatID:    chatID,
			Latitude:  loc.Lat,
			Longitude: loc.Long,
			Title:     loc.Title,
			Address:   loc.Description,
		})
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
