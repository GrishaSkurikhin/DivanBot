package bot

import (
	bothandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/bot-handlers"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	dialoghandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/dialog-handlers"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	mapstater "github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger/map-stater"
	inlinehandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/inline-handlers"
	inlinekeyboard "github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/inline-keyboard"

	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
)

type OperationsStorage interface {
	dialoghandlers.UserRegistrator
	bothandlers.FilmsRegsGetter
	bothandlers.PrevFilmsGetter
	bothandlers.FutureFilmsGetter
	dialoghandlers.UserDataChanger
	bothandlers.UserDataGetter
	dialoghandlers.FeedbackSender
	bothandlers.AboutInfoGetter
	bothandlers.IsUserRegChecker
}

type telegramBot struct {
	*bot.Bot
}

func New(token string, log logger.BotLogger, operationsStorage OperationsStorage) (*telegramBot, error) {
	st := mapstater.New()
	d := dialoger.New(st)
	d.AddDialog(dialoger.LeaveFeedbackDialog, dialoghandlers.LeaveFeedback(log, operationsStorage), 2)
	d.AddDialog(dialoger.RegDialog, dialoghandlers.RegUser(log, operationsStorage), 4)
	d.AddDialog(dialoger.ChangeDataDialog, dialoghandlers.ChangeData(log, operationsStorage), 2)

	opts := []bot.Option{
		bot.WithDefaultHandler(bothandlers.Default(log, d)),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, commands.Start, bot.MatchTypeExact, bothandlers.Start(log))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.RegUser, bot.MatchTypeExact, bothandlers.RegUserStart(log, d, operationsStorage))

	b.RegisterHandler(bot.HandlerTypeMessageText, commands.FutureFilms, bot.MatchTypeExact, bothandlers.FutureFilms(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuFutureFilms, bot.MatchTypeExact, bothandlers.FutureFilms(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.PrevFilms, bot.MatchTypeExact, bothandlers.PrevFilms(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuPrevFilms, bot.MatchTypeExact, bothandlers.PrevFilms(log, operationsStorage))

	b.RegisterHandler(bot.HandlerTypeMessageText, commands.ShowRegs, bot.MatchTypeExact, bothandlers.ShowRegs(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuShowRegs, bot.MatchTypeExact, bothandlers.ShowRegs(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.ShowData, bot.MatchTypeExact, bothandlers.ShowData(log, operationsStorage, d))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuShowData, bot.MatchTypeExact, bothandlers.ShowData(log, operationsStorage, d))

	b.RegisterHandler(bot.HandlerTypeMessageText, commands.LeaveFeedback, bot.MatchTypeExact, bothandlers.LeaveFeedbackStart(log, d, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuLeaveFeedback, bot.MatchTypeExact, bothandlers.LeaveFeedbackStart(log, d, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.About, bot.MatchTypeExact, bothandlers.About(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuAbout, bot.MatchTypeExact, bothandlers.About(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.Help, bot.MatchTypeExact, bothandlers.Help(log))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuHelp, bot.MatchTypeExact, bothandlers.Help(log))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataCancel, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.Cancel(log)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataName, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.ChangeDataStart(log, d)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataSurname, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.ChangeDataStart(log, d)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataGroup, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.ChangeDataStart(log, d)))

	return &telegramBot{b}, nil
}
