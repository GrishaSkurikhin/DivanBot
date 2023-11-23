package slider

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type OnSelect func(ctx context.Context, b *bot.Bot, query *models.CallbackQuery, slideID string)
type GetNewSlides func(ctx context.Context, b *bot.Bot, query *models.CallbackQuery) ([]Slide,error)

type Slide struct {
	ID    string
	Photo string
	Text  string
}

var (
	cmdPrev = "prev"
	cmdNext = "next"
)

type Slider struct {
	prefix      string
	slides      []Slide
	buttonsText []string
	buttonsData []string
}

func New(slides []Slide, prefix string, opts ...Option) *Slider {
	s := &Slider{
		prefix: prefix,
		slides: slides,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Slider) Show(ctx context.Context, b *bot.Bot, chatID any) (*models.Message, error) {
	slide := s.slides[0]

	sendParams := &bot.SendPhotoParams{
		ChatID:      chatID,
		Photo:       &models.InputFileString{Data: slide.Photo},
		Caption:     slide.Text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.buildKeyboard(),
	}

	return b.SendPhoto(ctx, sendParams)
}
