package slider

import (
	"strconv"

	"github.com/go-telegram/bot/models"
)

func (s *Slider) buildKeyboard() models.InlineKeyboardMarkup {
	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "\u00AB", CallbackData: s.prefix +cmdPrev},
				{Text: strconv.Itoa(1) + "/" + strconv.Itoa(len(s.slides)), CallbackData: s.slides[0].ID},
				{Text: "\u00BB", CallbackData: s.prefix +cmdNext},
			},
		},
	}

	var row []models.InlineKeyboardButton
	for i := range s.buttonsText {
		row = append(row, models.InlineKeyboardButton{Text: s.buttonsText[i], CallbackData: s.prefix + s.buttonsData[i]})
	}
	if len(row) > 0 {
		kb.InlineKeyboard = append(kb.InlineKeyboard, row)
	}

	return kb
}
