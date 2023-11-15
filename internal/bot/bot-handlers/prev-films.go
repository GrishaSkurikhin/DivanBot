package bothandlers

import (
	"context"
	"strconv"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	preparemessages "github.com/GrishaSkurikhin/DivanBot/internal/bot/prepare-messages"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type PrevFilmsGetter interface {
	GetPrevFims() ([]models.Film, error)
}

func PrevFilms(log logger.BotLogger, prevFilmsGetter PrevFilmsGetter) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler  = "PrevFilms"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
		)

		films, err := prevFilmsGetter.GetPrevFims()
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get previous films", err)
			return
		}
		slides := make([]slider.Slide, 0, len(films))
		for _, film := range films {
			slides = append(slides, slider.Slide{
				Text:  preparemessages.FilmDescriptionPrev(film),
				Photo: film.PosterURL,
			})
		}

		opts := []slider.Option{
			slider.OnCancel("Закрыть", true, sliderPrevOnCancel),
		}

		sl := slider.New(slides, opts...)
		_, err = sl.Show(ctx, b, strconv.Itoa(int(update.Message.Chat.ID)))
		if err != nil {
			messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}

func sliderPrevOnCancel(ctx context.Context, b *bot.Bot, message *botModels.Message) {}
