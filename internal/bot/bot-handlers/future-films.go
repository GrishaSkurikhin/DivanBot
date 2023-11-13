package bothandlers

import (
	"context"

	inlinehandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/inline-handlers"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	preparemessages "github.com/GrishaSkurikhin/DivanBot/internal/bot/prepare-messages"
	requestinfo "github.com/GrishaSkurikhin/DivanBot/internal/bot/request-info"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type FutureFilmsGetter interface {
	GetFutureFims() ([]models.Film, error)
	inlinehandlers.FilmRegistrator
}

func FutureFilms(log logger.BotLogger, futureFilmsGetter FutureFilmsGetter) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler            = "FutureFilms"
			username, inputMsg = requestinfo.Get(update)
			chatID             = int(update.Message.Chat.ID)
			userID             = int(update.Message.From.ID)
		)

		films, err := futureFilmsGetter.GetFutureFims()
		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get previous films", err)
			return
		}

		slides := make([]slider.Slide, 0, len(films))
		for _, film := range films {
			slides = append(slides, slider.Slide{
				Text:  preparemessages.FilmDescriptionFuture(film),
				Photo: film.PosterURL,
			})
		}

		opts := []slider.Option{
			slider.OnSelect("Записаться", false, inlinehandlers.RegOnFilm(log, futureFilmsGetter, films, userID)),
			slider.OnSelect("Место", false, inlinehandlers.ShowLocation(log, films)),
			slider.OnCancel("Закрыть", true, sliderFutureOnCancel),
		}

		sl := slider.New(slides, opts...)
		_, err = sl.Show(ctx, b, chatID)
		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to show slider", err)
			return
		}
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}

func sliderFutureOnCancel(ctx context.Context, b *bot.Bot, message *botModels.Message) {}

func Govno(ctx context.Context, b *bot.Bot, message *botModels.Message, item int) {}
