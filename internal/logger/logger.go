package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type StandartLogger interface {
	DEBUG(msg string)
	StdINFO(msg string)
	StdERROR(msg string, err error)
	StdFATAL(msg string, err error)
}

type BotLogger interface {
	DEBUG(msg string)
	BotINFO(handler string, username string, inputMsg string, infoMsg string)
	BotERROR(handler string, username string, inputMsg string, infoMsg string, err error)
}

type logger struct {
	zerolog.Logger
}

func New(appEnv string) logger {
	var log zerolog.Logger

	switch appEnv {
	case "dev":
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.Level(zerolog.DebugLevel)).
			With().
			Logger()
	case "prod":
		log = zerolog.New(os.Stderr).
			Level(zerolog.Level(zerolog.InfoLevel)).
			With().
			Timestamp().
			Logger()
	}
	return logger{log}
}

func (l logger) DEBUG(msg string) {
	l.Debug().Msg(msg)
}

func (l logger) StdINFO(msg string) {
	l.Info().Msg(msg)
}

func (l logger) StdERROR(msg string, err error) {
	l.Error().Err(err).Msg(msg)
}

func (l logger) StdFATAL(msg string, err error) {
	l.Fatal().Err(err).Msg(msg)
}

func (l logger) BotINFO(handler string, username string, inputMsg string, infoMsg string) {
	l.Info().
		Str("handler", handler).
		Str("username", username).
		Str("inputMsg", inputMsg).
		Msg(infoMsg)
}

func (l logger) BotERROR(handler string, username string, inputMsg string, infoMsg string, err error) {
	l.Error().
		Str("handler", handler).
		Str("username", username).
		Str("inputMsg", inputMsg).
		Err(err).
		Msg(infoMsg)
}
