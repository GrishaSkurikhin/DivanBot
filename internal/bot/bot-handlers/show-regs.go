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

type FilmsRegsGetter interface {
	GetFilmsRegs(userID uint64) ([]models.Film, error)
	inlinehandlers.FilmRegDeleter
	IsUserReg(userID uint64) (bool, error)
}

func ShowRegs(log logger.BotLogger, filmsRegsGetter FilmsRegsGetter) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler            = "ShowRegs"
			username, inputMsg = requestinfo.Get(update)
			userID             = uint64(update.Message.Chat.ID)
			chatID             = int(update.Message.Chat.ID)
		)


		isReg, err := filmsRegsGetter.IsUserReg(userID)
		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "failed to check user reg", err)
			return
		}

		if !isReg {
			messagesender.Info(ctx, b, update, log, handler, username, inputMsg, "Не зарегестрированы")
			log.BotINFO(handler, username, inputMsg, "successfully")
			return
		}

		films, err := filmsRegsGetter.GetFilmsRegs(userID)
		if err != nil {
			messagesender.Error(ctx, b, update, log, handler, username, inputMsg, "Ошибка")
			log.BotERROR(handler, username, inputMsg, "Failed to get previous films", err)
			return
		}

		if len(films) == 0 {
			messagesender.Info(ctx, b, update, log, handler, username, inputMsg, "Не записей на фильмы")
			log.BotINFO(handler, username, inputMsg, "successfully")
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
			slider.OnSelect("Отменить запись", true, inlinehandlers.CancelRegOnFilm(log, filmsRegsGetter, films, int(userID))),
			slider.OnSelect("Место", false, inlinehandlers.ShowLocation(log, films)),
			slider.OnCancel("Закрыть", true, sliderRegsOnCancel),
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

func sliderRegsOnCancel(ctx context.Context, b *bot.Bot, message *botModels.Message) {}
