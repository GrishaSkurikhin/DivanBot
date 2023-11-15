package keyboards

import (
	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	inlinekeyboard "github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/inline-keyboard"
	"github.com/GrishaSkurikhin/DivanBot/pkg/go-telegram/ui/keyboard"
	"github.com/go-telegram/bot"
)

func MainMenu() *keyboard.Keyboard {
	return keyboard.New().
		Row().
		Button(commands.MenuFutureFilms).
		Button(commands.MenuPrevFilms).
		Row().
		Button(commands.MenuShowRegs).
		Button(commands.MenuShowData).
		Row().
		Button(commands.MenuLeaveFeedback).
		Button(commands.MenuAbout).
		Button(commands.MenuHelp)
}

func DialogMenu() *keyboard.Keyboard {
	return keyboard.New().
		Row().
		Button(commands.Cancel)
}

func ChangeData(b *bot.Bot) *inlinekeyboard.Keyboard {
	return inlinekeyboard.New(b).
		Row().
		Button("Изменить имя", []byte(commands.ChangeDataName)).
		Button("Изменить фамилию", []byte(commands.ChangeDataSurname)).
		Button("Изменить группу", []byte(commands.ChangeDataGroup)).
		Row().
		Button("Закрыть", []byte(commands.ChangeDataCancel))
}
