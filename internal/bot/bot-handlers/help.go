package bothandlers

import (
	"context"

	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

const (
	helpInfo = `
Для записи на фильмы и оставления отзывов необходимо сначала зарегистрироваться с помощью команды /reg.

При регистрации необходимо указать свои имя, фамилию и группу.

Информацию "о себе" можно посмотреть и изменить с помощью команды /data.
Аккаунт после регистрации удалить нельзя.

Запись на фильмы происходит во вкладке "Что будем показывать".
Отмена записи во вкладке "Мои записи на фильмы".

Команды:
		/start - начало работы с ботом
		/reg - регистрация пользователя
		/films - планируемые фильмы
		/past - прошедшие фильмы
		/films_regs - список фильмов, на которые вы записаны
		/data - информация о вас
		/feedback - оставить отзыв
		/about - о киноклубе "Диван"
		/menu - вызвать главное меню
		/help - помощь
	`
)

func Help(log logger.BotLogger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *botModels.Update) {
		var (
			handler  = "Help"
			username = update.Message.From.Username
			inputMsg = update.Message.Text
			chatID   = update.Message.Chat.ID
		)

		messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, helpInfo)
		log.BotINFO(handler, username, inputMsg, "successfully")
	}
}
