package sliders

import (
	preparemessages "github.com/GrishaSkurikhin/DivanBot/internal/bot/prepare-messages"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"
)

func FutureFilms(films []models.Film) *slider.Slider {
	slides := make([]slider.Slide, 0, len(films))
	for _, film := range films {
		slides = append(slides, slider.Slide{
			Text:  preparemessages.FilmDescriptionFuture(film),
			Photo: film.PosterURL,
		})
	}

	opts := []slider.Option{
		slider.Button("Записаться", "reg"),
		slider.Button("Место", "location"),
		slider.Button("Закрыть", "cancel"),
	}
	return slider.New(slides, "future", opts...)
}

func PrevFilms(films []models.Film) *slider.Slider {
	slides := make([]slider.Slide, 0, len(films))
	for _, film := range films {
		slides = append(slides, slider.Slide{
			Text:  preparemessages.FilmDescriptionPrev(film),
			Photo: film.PosterURL,
		})
	}

	opts := []slider.Option{
		slider.Button("Закрыть", "cancel"),
	}

	return slider.New(slides, "prev", opts...)
}

func RegFilms(films []models.Film) *slider.Slider {
	slides := make([]slider.Slide, 0, len(films))
	for _, film := range films {
		slides = append(slides, slider.Slide{
			Text:  preparemessages.FilmDescriptionFuture(film),
			Photo: film.PosterURL,
		})
	}

	opts := []slider.Option{
		slider.Button("Отменить запись", "cancelreg"),
		slider.Button("Место", "location"),
		slider.Button("Закрыть", "cancel"),
	}

	return slider.New(slides, "userfilms", opts...)
}