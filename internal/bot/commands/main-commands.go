package commands

const (
	Start = "/start"
	Help = "/help"

	RegUser = "/reg"

	FutureFilms = "/films"
	PrevFilms = "/past"

	ShowRegs = "/films_regs"
	ShowData = "/data"

	LeaveFeedback = "/feedback"
	About = "/about"
)

func Commands() map[string]struct{} {
	return map[string]struct{} {
		Start: {},
		Help: {},
		RegUser: {},
		FutureFilms: {},
		PrevFilms: {},
		ShowRegs: {},
		ShowData: {},
		LeaveFeedback: {},
		About: {},
		MenuFutureFilms: {},
		MenuPrevFilms: {},
		MenuShowRegs: {},
		MenuShowData: {},
		MenuLeaveFeedback: {},
		MenuAbout: {},
		MenuHelp: {},
		ChangeDataName: {},
		ChangeDataSurname: {},
		ChangeDataGroup: {},
		FutureFilmsPrefix + FutureFilmsReg: {},
		FutureFilmsPrefix + FutureFilmsLocation: {},
		UserFilmsPrefix + UserFilmsCancelReg: {},
		UserFilmsPrefix + UserFilmsLocation: {},
	}
}