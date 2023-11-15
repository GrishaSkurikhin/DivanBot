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

func ShowLocation(log logger.BotLogger, films []models.Film) slider.OnSelect {
	return func(ctx context.Context, b *bot.Bot, query *botModels.CallbackQuery, item int) {
		var (
			handler  = "ShowLocation"
			username = query.Message.From.Username
			inputMsg = query.Message.Text
			chatID   = query.Message.Chat.ID
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
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
