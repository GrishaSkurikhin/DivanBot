package preparemessages

import (
	"fmt"
	"time"

	"github.com/GrishaSkurikhin/DivanBot/internal/models"
)

func FilmDescriptionPrev(film models.Film) string {
	return fmt.Sprintf("🎬 <strong>%s</strong>\n\n%s\n\n📅 %s\n📍 %s",
		film.Name,
		film.Description,
		getDate(film.ShowDate),
		film.Location.Title,
	)
}

func FilmDescriptionFuture(film models.Film) string {
	var isCloseInfo string
	if film.IsOpen {
		isCloseInfo = "🟢 Регистрация открыта"
	} else {
		isCloseInfo = "🔴 Регистрация закрыта"
	}

	return fmt.Sprintf("🎬 %s\n\n%s\n\n📅 %s\n📍 %s\n 👥 Всего мест: %d\n%s",
		film.Name,
		film.Description,
		getDate(film.ShowDate),
		film.Location.Title,
		film.PlacesNum,
		isCloseInfo,
	)
}

func getDate(dateTime time.Time) string {
	day := dateTime.Day()
	year := dateTime.Year()
	monthEng := dateTime.Month().String()
	var monthRu string
	switch monthEng {
	case "January":
		monthRu = "января"
	case "February":
		monthRu = "февраля"
	case "March":
		monthRu = "марта"
	case "April":
		monthRu = "апреля"
	case "May":
		monthRu = "мая"
	case "June":
		monthRu = "июня"
	case "July":
		monthRu = "июля"
	case "August":
		monthRu = "августа"
	case "September":
		monthRu = "сентября"
	case "October":
		monthRu = "октября"
	case "November":
		monthRu = "ноября"
	case "December":
		monthRu = "декабря"
	}
	return fmt.Sprintf("%d %s %d ", day, monthRu, year)
}
