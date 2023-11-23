package slider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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

func Callback(handler OnSelect, isDelete bool, getNewSlides GetNewSlides) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		newSlides, err := getNewSlides(ctx, b, update.CallbackQuery)
		if err != nil {
			fmt.Printf("error get new slides: %v", err)
			callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}

		if len(newSlides) == 0 {
			_, errDelete := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.ID,
			})
			if errDelete != nil {
				fmt.Printf("error delete message in callback, %v", errDelete)
				callbackAnswer(ctx, b, update.CallbackQuery)
				return
			}
			callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}

		item := int([]byte(update.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0][1].Text)[0] - '1')
		if item > len(newSlides)-1 {
			updateSlider(ctx, b, update.CallbackQuery, newSlides, 0, len(newSlides))
			callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}

		slideID := update.CallbackQuery.Message.ReplyMarkup.InlineKeyboard[0][1].CallbackData
		if slideID != newSlides[item].ID {
			updateSlider(ctx, b, update.CallbackQuery, newSlides, 0, len(newSlides))
			callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}

		reNext := regexp.MustCompile(".+" + regexp.QuoteMeta(cmdNext) + "$")
		if reNext.MatchString(update.CallbackQuery.Data) {
			if item == len(newSlides)-1 {
				updateSlider(ctx, b, update.CallbackQuery, newSlides, 0, len(newSlides))
			} else {
				updateSlider(ctx, b, update.CallbackQuery, newSlides, item+1, len(newSlides))
			}
			callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}

		rePrev := regexp.MustCompile(".+" + regexp.QuoteMeta(cmdPrev) + "$")
		if rePrev.MatchString(update.CallbackQuery.Data) {
			if item == 0 {
				updateSlider(ctx, b, update.CallbackQuery, newSlides, len(newSlides)-1, len(newSlides))
			} else {
				updateSlider(ctx, b, update.CallbackQuery, newSlides, item-1, len(newSlides))
			}
			callbackAnswer(ctx, b, update.CallbackQuery)
			return
		}

		if isDelete {
			_, errDelete := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.ID,
			})
			if errDelete != nil {
				fmt.Printf("error delete message in callback, %v", errDelete)
				callbackAnswer(ctx, b, update.CallbackQuery)
				return
			}
		}

		handler(ctx, b, update.CallbackQuery, slideID)
		callbackAnswer(ctx, b, update.CallbackQuery)
	}
}

func RegistrateCmdButtons(b *bot.Bot, prefix string, getNewSlides GetNewSlides) {
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, prefix+cmdPrev, bot.MatchTypeExact, Callback(nil, false, getNewSlides))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, prefix+cmdNext, bot.MatchTypeExact, Callback(nil, false, getNewSlides))
}

func updateSlider(ctx context.Context, b *bot.Bot, query *models.CallbackQuery, slides []Slide, currSlide int, slidesCount int) {
	reply := query.Message.ReplyMarkup
	reply.InlineKeyboard[0][1].Text = fmt.Sprintf("%d/%d", currSlide+1, slidesCount)
	reply.InlineKeyboard[0][1].CallbackData = slides[currSlide].ID

	editParams := &bot.EditMessageMediaParams{
		ChatID:    query.Message.Chat.ID,
		MessageID: query.Message.ID,
		Media: &models.InputMediaPhoto{
			Media:     slides[currSlide].Photo,
			Caption:   slides[currSlide].Text,
			ParseMode: models.ParseModeHTML,
		},
		ReplyMarkup: reply,
	}

	_, errEdit := b.EditMessageMedia(ctx, editParams)
	if errEdit != nil {
		if strings.Contains(errEdit.Error(), "message is not modified") {
			return
		}
		fmt.Printf("error edit message in callback, %v", errEdit)
	}
}
