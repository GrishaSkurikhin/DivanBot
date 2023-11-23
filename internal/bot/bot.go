package bot

import (
	"context"

	bothandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/bot-handlers"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	dialoghandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/dialog-handlers"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/dialoger"
	inlinehandlers "github.com/GrishaSkurikhin/DivanBot/internal/bot/inline-handlers"
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/middlewares"
	inlinekeyboard "github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/inline-keyboard"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/slider"

	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
	inlinehandlers.LocationGetter
}

type telegramBot struct {
	*bot.Bot
}

func New(token string, log logger.BotLogger, operationsStorage OperationsStorage, st dialoger.Stater) (*telegramBot, error) {
	d := dialoger.New(st)
	d.AddDialog(dialoger.LeaveFeedbackDialog, dialoghandlers.LeaveFeedback(log, operationsStorage), 2)
	d.AddDialog(dialoger.RegDialog, dialoghandlers.RegUser(log, operationsStorage), 5)
	d.AddDialog(dialoger.ChangeDataDialog, dialoghandlers.ChangeData(log, operationsStorage), 2)

	opts := []bot.Option{
		bot.WithMiddlewares(middlewares.CheckDialog(log, d)),
		bot.WithDefaultHandler(bothandlers.Default(log, d)),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	registrateStandartHandlers(b, log, operationsStorage, d)
	registrateMenuHandlers(b, log, operationsStorage, d)

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataCancel, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.Cancel(log)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataName, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.ChangeDataStart(log, d)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataSurname, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.ChangeDataStart(log, d)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.ChangeDataGroup, bot.MatchTypeContains, inlinekeyboard.Callback(inlinehandlers.ChangeDataStart(log, d)))

	slider.RegistrateCmdButtons(b, commands.PrevFilmsPrefix, inlinehandlers.GetPrevFilmsSlides(operationsStorage))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.PrevFilmsPrefix+commands.PrevFilmsCancel, bot.MatchTypeExact, slider.Callback(func(ctx context.Context, b *bot.Bot, query *models.CallbackQuery, slideID string) {}, true, inlinehandlers.GetPrevFilmsSlides(operationsStorage)))

	slider.RegistrateCmdButtons(b, commands.FutureFilmsPrefix, inlinehandlers.GetFutureFilmsSlides(operationsStorage))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.FutureFilmsPrefix+commands.FutureFilmsCancel, bot.MatchTypeExact, slider.Callback(func(ctx context.Context, b *bot.Bot, query *models.CallbackQuery, slideID string) {}, true, inlinehandlers.GetFutureFilmsSlides(operationsStorage)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.FutureFilmsPrefix+commands.FutureFilmsReg, bot.MatchTypeExact, slider.Callback(inlinehandlers.RegOnFilm(log, operationsStorage), false, inlinehandlers.GetFutureFilmsSlides(operationsStorage)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.FutureFilmsPrefix+commands.FutureFilmsLocation, bot.MatchTypeExact, slider.Callback(inlinehandlers.ShowLocation(log, operationsStorage), false, inlinehandlers.GetFutureFilmsSlides(operationsStorage)))

	slider.RegistrateCmdButtons(b, commands.UserFilmsPrefix, inlinehandlers.GetUserFilmsSlides(operationsStorage))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.UserFilmsPrefix+commands.UserFilmsCancel, bot.MatchTypeExact, slider.Callback(func(ctx context.Context, b *bot.Bot, query *models.CallbackQuery, slideID string) {}, true, inlinehandlers.GetUserFilmsSlides(operationsStorage)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.UserFilmsPrefix+commands.UserFilmsCancelReg, bot.MatchTypeExact, slider.Callback(inlinehandlers.CancelRegOnFilm(log, operationsStorage), true, inlinehandlers.GetUserFilmsSlides(operationsStorage)))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, commands.UserFilmsPrefix+commands.UserFilmsLocation, bot.MatchTypeExact, slider.Callback(inlinehandlers.ShowLocation(log, operationsStorage), false, inlinehandlers.GetUserFilmsSlides(operationsStorage)))
	return &telegramBot{b}, nil
}

func registrateStandartHandlers(b *bot.Bot, log logger.BotLogger, operationsStorage OperationsStorage, d *dialoger.Dialoger) {
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.Start, bot.MatchTypeExact, bothandlers.Start(log))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.RegUser, bot.MatchTypeExact, bothandlers.RegUserStart(log, d, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.FutureFilms, bot.MatchTypeExact, bothandlers.FutureFilms(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.PrevFilms, bot.MatchTypeExact, bothandlers.PrevFilms(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.ShowRegs, bot.MatchTypeExact, bothandlers.ShowRegs(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.ShowData, bot.MatchTypeExact, bothandlers.ShowData(log, operationsStorage, d))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.LeaveFeedback, bot.MatchTypeExact, bothandlers.LeaveFeedbackStart(log, d, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.About, bot.MatchTypeExact, bothandlers.About(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.Help, bot.MatchTypeExact, bothandlers.Help(log))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MainMenu, bot.MatchTypeExact, bothandlers.MainMenu(log))
}

func registrateMenuHandlers(b *bot.Bot, log logger.BotLogger, operationsStorage OperationsStorage, d *dialoger.Dialoger) {
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuFutureFilms, bot.MatchTypeExact, bothandlers.FutureFilms(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuShowRegs, bot.MatchTypeExact, bothandlers.ShowRegs(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuPrevFilms, bot.MatchTypeExact, bothandlers.PrevFilms(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuShowData, bot.MatchTypeExact, bothandlers.ShowData(log, operationsStorage, d))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuLeaveFeedback, bot.MatchTypeExact, bothandlers.LeaveFeedbackStart(log, d, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuAbout, bot.MatchTypeExact, bothandlers.About(log, operationsStorage))
	b.RegisterHandler(bot.HandlerTypeMessageText, commands.MenuHelp, bot.MatchTypeExact, bothandlers.Help(log))
}
