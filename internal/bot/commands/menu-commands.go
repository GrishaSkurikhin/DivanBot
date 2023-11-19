package commands

const (
	MenuFutureFilms = "Что будем показывать"
	MenuPrevFilms   = "Что уже показали"

	MenuShowRegs = "Мои записи на фильмы"
	MenuShowData = "Мои данные"

	MenuLeaveFeedback = "Оставить отзыв"
	MenuAbout         = `О "divan"`
	MenuHelp          = "Помощь с ботом"

	Cancel  = "Отмена"
)

func MenuCommands() map[string]struct{} {
	return map[string]struct{} {
		MenuFutureFilms: {},
		MenuPrevFilms: {},
		MenuShowRegs: {},
		MenuShowData: {},
		MenuLeaveFeedback: {},
		MenuAbout: {},
		MenuHelp: {},
	}
}