package keyboard

import (
	"encoding/json"

	"github.com/go-telegram/bot/models"
)

type OnErrorHandler func(err error)

type Keyboard struct {
	markup [][]models.KeyboardButton
}

func New() *Keyboard {
	kb := &Keyboard{
		markup: [][]models.KeyboardButton{{}},
	}

	return kb
}

func (kb *Keyboard) MarshalJSON() ([]byte, error) {
	return json.Marshal(models.ReplyKeyboardMarkup{
		Keyboard:       kb.markup,
		ResizeKeyboard: true,
	})
}

func (kb *Keyboard) Row() *Keyboard {
	if len(kb.markup[len(kb.markup)-1]) > 0 {
		kb.markup = append(kb.markup, []models.KeyboardButton{})
	}
	return kb
}

func (kb *Keyboard) Button(text string) *Keyboard {
	kb.markup[len(kb.markup)-1] = append(kb.markup[len(kb.markup)-1], models.KeyboardButton{
		Text: text,
	})

	return kb
}
