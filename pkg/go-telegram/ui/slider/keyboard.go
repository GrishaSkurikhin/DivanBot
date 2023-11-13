package slider

import (
	"strconv"

	"github.com/go-telegram/bot/models"
)

func (s *Slider) buildKeyboard() models.InlineKeyboardMarkup {
	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "\u00AB", CallbackData: s.prefix + cmdPrev},
				{Text: strconv.Itoa(s.current+1) + "/" + strconv.Itoa(len(s.slides)), CallbackData: s.prefix + cmdNop},
				{Text: "\u00BB", CallbackData: s.prefix + cmdNext},
			},
		},
	}

	var row []models.InlineKeyboardButton
	for i := range s.onSelect {
		row = append(row, models.InlineKeyboardButton{Text: s.selectButtonsText[i], CallbackData: s.prefix + cmdSelect + strconv.Itoa(i)})
	}
	if s.onCancel != nil {
		row = append(row, models.InlineKeyboardButton{Text: s.cancelButtonText, CallbackData: s.prefix + cmdCancel})
	}

	if len(row) > 0 {
		kb.InlineKeyboard = append(kb.InlineKeyboard, row)
	}

	return kb
}