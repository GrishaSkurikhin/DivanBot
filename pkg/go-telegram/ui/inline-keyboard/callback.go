package inlinekeyboard

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
		fmt.Println("callback answer failed")
	}
}

func Callback(handler OnSelect) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		_, errDelete := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.ID,
		})
		if errDelete != nil {
			fmt.Printf("error delete message in callback, %v", errDelete)
			callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}

		handler(ctx, b, update.CallbackQuery)
		callbackAnswer(ctx, b, update.CallbackQuery)
	}
}
