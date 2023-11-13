package requestinfo

import "github.com/go-telegram/bot/models"

func Get(update *models.Update) (username string, inputTxt string) {
	username = update.Message.From.Username
	inputTxt = update.Message.Text
	return
}
