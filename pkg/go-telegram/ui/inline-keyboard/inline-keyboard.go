package inlinekeyboard

import (
	"context"
	"encoding/json"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type OnSelect func(ctx context.Context, bot *bot.Bot, mes *models.Message, data []byte)

type Keyboard struct {
	prefix            string
	data              [][]byte
	callbackHandlerID string
	markup            [][]models.InlineKeyboardButton
}

func New(b *bot.Bot, prefix string) *Keyboard {
	return &Keyboard{
		prefix: prefix,
		markup: [][]models.InlineKeyboardButton{{}},
		data:   [][]byte{},
	}
}

func (kb *Keyboard) MarshalJSON() ([]byte, error) {
	return json.Marshal(models.InlineKeyboardMarkup{InlineKeyboard: kb.markup})
}
