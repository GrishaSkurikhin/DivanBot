package inlinehandlers

import (
	"context"

	preparemessages "github.com/GrishaSkurikhin/DivanBot/internal/bot/prepare-messages"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type PrevFilmsGetter interface {
	GetPrevFims() ([]models.Film, error)
}

type FutureFilmsGetter interface {
	GetFutureFims() ([]models.Film, error)
}

type FilmsRegsGetter interface {
	GetFilmsRegs(userID uint64) ([]models.Film, error)
}

func GetPrevFilmsSlides(prevFilmsGetter PrevFilmsGetter) slider.GetNewSlides {
	return func(ctx context.Context, b *bot.Bot, query *botModels.CallbackQuery) ([]slider.Slide, error) {
		films, err := prevFilmsGetter.GetPrevFims()
		if err != nil {
			return nil, err
		}

		slides := make([]slider.Slide, 0, len(films))
		for _, film := range films {
			slides = append(slides, slider.Slide{
				ID: film.ID,
				Text:  preparemessages.FilmDescriptionPrev(film),
				Photo: film.PosterURL,
			})
		}
		return slides,  nil
	}
}

func GetFutureFilmsSlides(futureFilmsGetter FutureFilmsGetter) slider.GetNewSlides {
	return func(ctx context.Context, b *bot.Bot, query *botModels.CallbackQuery) ([]slider.Slide, error) {
		films, err := futureFilmsGetter.GetFutureFims()
		if err != nil {
			return nil, err
		}

		slides := make([]slider.Slide, 0, len(films))
		for _, film := range films {
			slides = append(slides, slider.Slide{
				ID: film.ID,
				Text:  preparemessages.FilmDescriptionFuture(film),
				Photo: film.PosterURL,
			})
		}
		return slides, nil
	}
}

func GetUserFilmsSlides(filmsRegsGetter FilmsRegsGetter) slider.GetNewSlides {
	return func(ctx context.Context, b *bot.Bot, query *botModels.CallbackQuery) ([]slider.Slide, error) {
		userID := uint64(query.Message.Chat.ID)
		films, err := filmsRegsGetter.GetFilmsRegs(userID)
		if err != nil {
			return nil, err
		}

		slides := make([]slider.Slide, 0, len(films))
		for _, film := range films {
			slides = append(slides, slider.Slide{
				ID: film.ID,
				Text:  preparemessages.FilmDescriptionFuture(film),
				Photo: film.PosterURL,
			})
		}
		return slides, nil
	}
}