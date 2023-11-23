package dialoghandlers

import (
	"context"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/keyboards"
	messagesender "github.com/GrishaSkurikhin/DivanBot/internal/bot/message-sender"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/go-telegram/bot"
	botModels "github.com/go-telegram/bot/models"
)

type UserRegistrator interface {
	RegUser(models.User) error
}

func RegUser(log logger.BotLogger, userRegistrator UserRegistrator) dialoger.DialogHandler {
	return func(ctx context.Context, b *bot.Bot, msg *botModels.Message, state int, info map[string]string) (newInfo map[string]string, isErr bool) {
		var (
			handler  = "RegUser"
			username = msg.From.Username
			inputMsg = msg.Text
			userID = uint64(msg.From.ID)
			chatID = msg.Chat.ID
		)
		newInfo = make(map[string]string)

		switch state {
		case 1:
			messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username,
				inputMsg, "Введите имя", keyboards.DialogMenu())
		case 2:
			newInfo["name"] = inputMsg
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Введите фамилию")
		case 3:
			newInfo["surname"] = inputMsg
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, "Введите группу")
		case 4:
			newInfo["group"] = inputMsg
			messagesender.Info(ctx, b, chatID, log, handler, username, inputMsg, `Осталось еще чуть-чуть. Введите, пожалуйста, откуда вы про нас узнали. Если не хотите отвечать, просто отправьте "-"`)
		case 5:
			user := models.User{
				TgID:     userID,
				Username: username,
				Name:     info["name"],
				Surname:  info["surname"],
				Group: info["group"],    
				WhereFind: inputMsg,
			}
			err := userRegistrator.RegUser(user)
			if err != nil {
				messagesender.Error(ctx, b, chatID, log, handler, username, inputMsg, "Ошибка")
				log.BotERROR(handler, username, inputMsg, "Failed to send user to db", err)
				isErr = true
				return
			}
			messagesender.InfoWithKeyboard(ctx, b, chatID, log, handler, username, inputMsg,
				"Вы успешно зарегестрированы", keyboards.MainMenu())
		}
	

		log.BotINFO(handler, username, inputMsg, "successfully")
		return
	}
}
