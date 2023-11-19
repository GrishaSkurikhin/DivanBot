package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot"
	"github.com/GrishaSkurikhin/DivanBot/internal/config"
	"github.com/GrishaSkurikhin/DivanBot/internal/logger"
	"github.com/GrishaSkurikhin/DivanBot/internal/storage/ydb"
)

const (
	shutdownTimeout = 3 * time.Second
)

func main() {
	conf := config.New()

	log := logger.New("dev")
	log.DEBUG("debug messages are enabled")
	log.StdINFO("starting bot")

	ydb, err := ydb.NewWithServiceAccount(conf.YDB.DSN, conf.YDB.KeyPath)
	if err != nil {
		log.StdFATAL("failed to connect ydb", err)
		os.Exit(1)
	}

	b, err := bot.New(conf.Token, log, ydb, ydb)
	if err != nil {
		log.StdFATAL("failed to create bot", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go b.Start(ctx)
	log.StdINFO("bot started")

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err = ydb.Close(shutdownCtx)
	if err != nil {
		log.StdFATAL("error while close db connection", err)
		os.Exit(1)
	}

	log.StdINFO("bot stopped successfully")
}
