package ydb

import (
	"context"
	"fmt"
	"time"

	customerrors "github.com/GrishaSkurikhin/DivanBot/internal/custom-errors"
	"github.com/GrishaSkurikhin/DivanBot/internal/models"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc"
)

const (
	connectionDeadline = 20 * time.Second
	operationDeadline  = 20 * time.Second
)

type yandexDatabase struct {
	conn *ydb.Driver
}

func NewWithToken(dsn string, IAMtoken string) (*yandexDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionDeadline)
	defer cancel()

	db, err := ydb.Open(ctx, dsn,
		ydb.WithAccessTokenCredentials(IAMtoken),
	)
	if err != nil {
		return nil, err
	}
	return &yandexDatabase{conn: db}, nil
}

func NewWithServiceAccount(dsn string, filePath string) (*yandexDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionDeadline)
	defer cancel()

	db, err := ydb.Open(ctx, dsn,
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials(filePath),
	)
	if err != nil {
		return nil, err
	}
	return &yandexDatabase{conn: db}, nil
}

func (ydb *yandexDatabase) Close(ctx context.Context) error {
	return ydb.conn.Close(ctx)
}

func (ydb *yandexDatabase) RegUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $UserID AS Uint64;
	DECLARE $Username AS String;
	DECLARE $Name AS Utf8;
	DECLARE $Surname AS Utf8;
	DECLARE $Group AS Utf8;
	DECLARE $WhereFind AS Utf8;
	UPSERT INTO Users ( tg_id, username, name, surname, group, where_find )
	VALUES ( $UserID, $Username, $Name, $Surname, $Group, $WhereFind );`

	params := table.NewQueryParameters(
		table.ValueParam("$UserID", types.Uint64Value(user.TgID)),
		table.ValueParam("$Username", types.BytesValueFromString(user.Username)),
		table.ValueParam("$Name", types.UTF8Value(user.Name)),
		table.ValueParam("$Surname", types.UTF8Value(user.Surname)),
		table.ValueParam("$Group", types.UTF8Value(user.Group)),
		table.ValueParam("$WhereFind", types.UTF8Value(user.WhereFind)),
	)

	return ydb.CUDoperations(ctx, query, params)
}

func (ydb *yandexDatabase) IsUserReg(userID uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $userID AS Uint64;
	SELECT EXISTS(SELECT tg_id FROM Users WHERE tg_id = $userID)`

	params := table.NewQueryParameters(
		table.ValueParam("$userID", types.Uint64Value(userID)),
	)

	info, err := ydb.selectSingleValue(ctx, query, params)
	return info.(bool), err
}

func (ydb *yandexDatabase) GetUserData(userID uint64) (string, string, string, error) {
	var (
		name    *string
		surname *string
		group   *string
	)

	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $UserID AS Uint64;
	SELECT u.name, u.surname, u.group
	FROM Users AS u
	WHERE u.tg_id = $UserID`

	params := table.NewQueryParameters(
		table.ValueParam("$UserID", types.Uint64Value(userID)),
	)

	readTx := table.TxControl(table.BeginTx(table.WithOnlineReadOnly()), table.CommitTx())
	err := ydb.conn.Table().Do(ctx,
		func(ctx context.Context, s table.Session) error {
			_, res, err := s.Execute(ctx, readTx, query, params)
			if err != nil {
				return fmt.Errorf("execute error: %v", err)
			}
			defer res.Close()

			if res.NextResultSet(ctx) {
				if res.NextRow() {
					err = res.ScanNamed(
						named.Optional("name", &name),
						named.Optional("surname", &surname),
						named.Optional("group", &group),
					)
					if err != nil {
						return fmt.Errorf("scan error: %v", err)
					}
				}
			}
			
			return res.Err()
		},
	)

	if name == nil {
		return "", "", "", customerrors.UserNotRegistered{}
	}
	return *name, *surname, *group, err
}

func (ydb *yandexDatabase) ChangeUserData(dataType string, newValue string, userID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := fmt.Sprintf(`DECLARE $UserID AS Uint64;
	DECLARE $Value AS Utf8;
	UPDATE Users
	SET %s=$Value
	WHERE tg_id = $UserID;`, dataType)

	params := table.NewQueryParameters(
		table.ValueParam("$UserID", types.Uint64Value(userID)),
		table.ValueParam("$Value", types.UTF8Value(newValue)),
	)

	return ydb.CUDoperations(ctx, query, params)
}

func (ydb *yandexDatabase) IsExistRegOnFilm(userID uint64, filmID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $userID AS Uint64;
	DECLARE $filmID AS String;
	SELECT EXISTS(SELECT user_tg_id FROM Registrations 
	WHERE user_tg_id = $userID AND film_id =  $filmID)`

	params := table.NewQueryParameters(
		table.ValueParam("$userID", types.Uint64Value(userID)),
		table.ValueParam("$filmID", types.BytesValueFromString(filmID)),
	)

	info, err := ydb.selectSingleValue(ctx, query, params)
	return info.(bool), err
}

func (ydb *yandexDatabase) RegOnFilm(userID uint64, filmID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $UserID AS Uint64;
	DECLARE $FilmID AS String;
	UPSERT INTO Registrations ( user_tg_id, film_id )
	VALUES ( $UserID, $FilmID );`

	params := table.NewQueryParameters(
		table.ValueParam("$UserID", types.Uint64Value(userID)),
		table.ValueParam("$FilmID", types.BytesValueFromString(filmID)),
	)

	return ydb.CUDoperations(ctx, query, params)
}

func (ydb *yandexDatabase) DeleteRegOnFilm(userID uint64, filmID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $UserID AS Uint64;
	DECLARE $FilmID AS String;
	DELETE FROM Registrations
	WHERE user_tg_id = $UserID AND film_id = $FilmID;`

	params := table.NewQueryParameters(
		table.ValueParam("$UserID", types.Uint64Value(userID)),
		table.ValueParam("$FilmID", types.BytesValueFromString(filmID)),
	)

	return ydb.CUDoperations(ctx, query, params)
}

func (ydb *yandexDatabase) GetLocation(filmID string) (models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $FilmID AS String;
	SELECT
	l.title as title,
	l.description as description,
	l.latitude as latitude,
	l.longitude as longitude,
	l.video_url as video_url
	FROM Locations AS l
	JOIN Films AS f ON f.location_id = l.id
	WHERE f.id = $FilmID`

	params := table.NewQueryParameters(
		table.ValueParam("$FilmID", types.BytesValueFromString(filmID)),
	)

	return ydb.selectLocation(ctx, query, params)
}

func (ydb *yandexDatabase) GetFilmsRegs(userID uint64) ([]models.Film, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $UserID AS Uint64;
	SELECT
	f.id as film_id,
	f.name as film_name,
	f.description as film_description,
	f.show_date as film_show_date,
	f.poster_url as film_poster_url,
	f.is_open as film_is_open,
	f.places_num as film_places_num,
	l.title as location_title,
	l.description as location_description,
	l.latitude as location_latitude,
	l.longitude as location_longitude,
	l.video_url as location_video_url
	FROM Films AS f
	JOIN Locations AS l ON f.location_id = l.id
	JOIN Registrations AS r ON f.id = r.film_id
	WHERE r.user_tg_id = $UserID
	ORDER BY film_show_date`

	params := table.NewQueryParameters(
		table.ValueParam("$UserID", types.Uint64Value(userID)),
	)

	return ydb.selectFilms(ctx, query, params)
}

func (ydb *yandexDatabase) GetPrevFims() ([]models.Film, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `SELECT
	f.id as film_id,
	f.name as film_name,
	f.description as film_description,
	f.show_date as film_show_date,
	f.poster_url as film_poster_url,
	f.is_open as film_is_open,
	f.places_num as film_places_num,
	l.title as location_title,
	l.description as location_description,
	l.latitude as location_latitude,
	l.longitude as location_longitude,
	l.video_url as location_video_url
	FROM Films AS f
	JOIN Locations AS l ON f.location_id = l.id
	WHERE f.show_date < CurrentUtcDatetime()
	ORDER BY film_show_date`

	return ydb.selectFilms(ctx, query, nil)
}

func (ydb *yandexDatabase) GetFutureFims() ([]models.Film, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `SELECT
	f.id as film_id,
	f.name as film_name,
	f.description as film_description,
	f.show_date as film_show_date,
	f.poster_url as film_poster_url,
	f.is_open as film_is_open,
	f.places_num as film_places_num,
	l.title as location_title,
	l.description as location_description,
	l.latitude as location_latitude,
	l.longitude as location_longitude,
	l.video_url as location_video_url
	FROM Films AS f
	JOIN Locations AS l ON f.location_id = l.id
	WHERE f.show_date > CurrentUtcDatetime()
	ORDER BY film_show_date`

	return ydb.selectFilms(ctx, query, nil)
}

func (ydb *yandexDatabase) SendFeedback(userID uint64, comment string) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $UUID AS String;
	DECLARE $UserID AS Uint64;
	DECLARE $Comment AS Utf8;
	UPSERT INTO Feedbacks ( uuid, user_tg_id, comment )
	VALUES ( $UUID, $UserID, $Comment );`

	params := table.NewQueryParameters(
		table.ValueParam("UUID", types.BytesValueFromString(uuid.New().String())),
		table.ValueParam("$UserID", types.Uint64Value(userID)),
		table.ValueParam("$Comment", types.UTF8Value(comment)),
	)

	return ydb.CUDoperations(ctx, query, params)
}

func (ydb *yandexDatabase) GetAboutInfo() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `SELECT about FROM Info`
	info, err := ydb.selectSingleValue(ctx, query, nil)
	return info.(string), err

}

func (ydb *yandexDatabase) CUDoperations(ctx context.Context, query string, params *table.QueryParameters) error {
	return ydb.conn.Table().DoTx(ctx,
		func(ctx context.Context, tx table.TransactionActor) error {
			res, err := tx.Execute(ctx, query, params)
			if err != nil {
				return fmt.Errorf("execute error: %v", err)
			}
			defer res.Close()

			return res.Err()
		}, table.WithIdempotent(),
	)
}

func (ydb *yandexDatabase) selectSingleValue(ctx context.Context, query string, params *table.QueryParameters) (any, error) {
	var value any
	readTx := table.TxControl(table.BeginTx(table.WithOnlineReadOnly()), table.CommitTx())

	err := ydb.conn.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			_, res, err := s.Execute(ctx, readTx, query, params)
			if err != nil {
				return fmt.Errorf("execute error: %v", err)
			}
			defer res.Close()

			if res.NextResultSet(ctx) {
				if res.NextRow() {
					err = res.ScanWithDefaults(&value)
					if err != nil {
						return fmt.Errorf("scan res error: %v", err)
					}
				}
			}
			return res.Err()
		},
	)
	return value, err
}

func (ydb *yandexDatabase) selectLocation(ctx context.Context, query string, params *table.QueryParameters) (models.Location, error) {
	location := models.Location{}

	readTx := table.TxControl(table.BeginTx(table.WithOnlineReadOnly()), table.CommitTx())
	err := ydb.conn.Table().Do(ctx,
		func(ctx context.Context, s table.Session) error {
			var (
				title       *string
				description *string
				latitude    float64
				longitude   float64
				video_url   *string
			)
			_, res, err := s.Execute(ctx, readTx, query, params)
			if err != nil {
				return fmt.Errorf("execute error: %v", err)
			}
			defer res.Close()

			if res.NextResultSet(ctx) {
				if res.NextRow() {
					err = res.ScanNamed(
						named.Optional("title", &title),
						named.Optional("description", &description),
						named.Required("latitude", &latitude),
						named.Required("longitude", &longitude),
						named.Optional("video_url", &video_url),
					)
					if err != nil {
						return fmt.Errorf("scan error: %v", err)
					}
				}
			}
			if res.Err() != nil {
				return res.Err()
			}
			location.Title = *title
			location.Description = *description
			location.Lat = latitude
			location.Long = longitude
			location.VideoURL = *video_url
			return nil
		},
	)
	return location, err
}

func (ydb *yandexDatabase) selectFilms(ctx context.Context, query string, params *table.QueryParameters) ([]models.Film, error) {
	films := []models.Film{}

	readTx := table.TxControl(table.BeginTx(table.WithOnlineReadOnly()), table.CommitTx())
	err := ydb.conn.Table().Do(ctx,
		func(ctx context.Context, s table.Session) error {
			_, res, err := s.Execute(ctx, readTx, query, params)
			if err != nil {
				return fmt.Errorf("execute error: %v", err)
			}
			defer res.Close()

			for res.NextResultSet(ctx) {
				for res.NextRow() {
					film, err := scanFilm(res)
					if err != nil {
						return fmt.Errorf("scan error: %v", err)
					}
					films = append(films, film)
				}
			}
			return res.Err()
		},
	)
	return films, err
}

func scanFilm(res result.Result) (models.Film, error) {
	var (
		film_id              *string
		film_name            *string
		film_description     *string
		film_show_date       *time.Time
		film_poster_url      *string
		film_is_open         bool
		film_places_num      uint64
		location_title       *string
		location_description *string
		location_latitude    float64
		location_longitude   float64
		location_video_url   *string
	)

	err := res.ScanNamed(
		named.Optional("film_id", &film_id),
		named.Optional("film_name", &film_name),
		named.Optional("film_description", &film_description),
		named.Optional("film_show_date", &film_show_date),
		named.Optional("film_poster_url", &film_poster_url),
		named.Required("film_is_open", &film_is_open),
		named.Required("film_places_num", &film_places_num),
		named.Optional("location_title", &location_title),
		named.Optional("location_description", &location_description),
		named.Required("location_latitude", &location_latitude),
		named.Required("location_longitude", &location_longitude),
		named.Optional("location_video_url", &location_video_url),
	)
	if err != nil {
		return models.Film{}, err
	}

	return models.Film{
		ID:          *film_id,
		Name:        *film_name,
		Description: *film_description,
		ShowDate:    *film_show_date,
		PosterURL:   *film_poster_url,
		IsOpen:      film_is_open,
		PlacesNum:   film_places_num,
		Location: models.Location{
			Title:       *location_title,
			Description: *location_description,
			Lat:         location_latitude,
			Long:        location_longitude,
			VideoURL:    *location_video_url,
		},
	}, nil
}
