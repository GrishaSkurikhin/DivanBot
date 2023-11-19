package sliders

import (
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	preparemessages "github.com/GrishaSkurikhin/DivanBot/internal/bot/prepare-messages"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
)

func FutureFilms(films []models.Film) *slider.Slider {
	slides := make([]slider.Slide, 0, len(films))
	for _, film := range films {
		slides = append(slides, slider.Slide{
			ID:    film.ID,
			Text:  preparemessages.FilmDescriptionFuture(film),
			Photo: film.PosterURL,
		})
	}

	opts := []slider.Option{
		slider.Button("Записаться", commands.FutureFilmsReg),
		slider.Button("Место", commands.FutureFilmsLocation),
		slider.Button("Закрыть", commands.FutureFilmsCancel),
	}
	return slider.New(slides, commands.FutureFilmsPrefix, opts...)
}

func PrevFilms(films []models.Film) *slider.Slider {
	slides := make([]slider.Slide, 0, len(films))
	for _, film := range films {
		slides = append(slides, slider.Slide{
			ID:    film.ID,
			Text:  preparemessages.FilmDescriptionPrev(film),
			Photo: film.PosterURL,
		})
	}

	opts := []slider.Option{
		slider.Button("Закрыть", commands.PrevFilmsCancel),
	}

	return slider.New(slides, commands.PrevFilmsPrefix, opts...)
}

func RegFilms(films []models.Film) *slider.Slider {
	slides := make([]slider.Slide, 0, len(films))
	for _, film := range films {
		slides = append(slides, slider.Slide{
			ID:    film.ID,
			Text:  preparemessages.FilmDescriptionFuture(film),
			Photo: film.PosterURL,
		})
	}

	opts := []slider.Option{
		slider.Button("Отменить запись", commands.UserFilmsCancelReg),
		slider.Button("Место", commands.UserFilmsLocation),
		slider.Button("Закрыть", commands.UserFilmsCancel),
	}

	return slider.New(slides, commands.UserFilmsPrefix, opts...)
}
