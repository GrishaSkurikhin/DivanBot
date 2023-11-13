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
		Button(commands.FutureFilms).
		Button(commands.PrevFilms).
		Row().
		Button(commands.ShowRegs).
		Button(commands.ShowData).
		Row().
		Button(commands.LeaveFeedback).
		Button(commands.About).
		Button(commands.Help)
}

func DialogMenu() *keyboard.Keyboard {
	return keyboard.New().
		Row().
		Button(commands.Cancel)
}

func ChangeData(b *bot.Bot) *inlinekeyboard.Keyboard {
	return inlinekeyboard.New(b, "data").
		Row().
		Button("Изменить имя", []byte("name")).
		Button("Изменить фамилию", []byte("surname")).
		Button("Изменить группу", []byte("group")).
		Row().
		Button("Закрыть", []byte("close"))
}
