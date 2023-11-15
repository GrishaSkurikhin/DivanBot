package config

import (
	"os"
)

type bot struct {
	Token string
	YDB   ydb
}

type ydb struct {
	DSN   string
	KeyPath string
}

func New() bot {
	return bot{
		Token: os.Getenv("TG_TOKEN"),
		YDB: ydb{
			DSN:   os.Getenv("YDB_DSN"),
			KeyPath: os.Getenv("YDB_KEY_PATH"),
		},
	}
}
