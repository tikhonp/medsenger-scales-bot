package main

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	medsengerscalesbot "github.com/tikhonp/medsenger-scales-bot"
	"github.com/tikhonp/medsenger-scales-bot/config"
	"github.com/tikhonp/medsenger-scales-bot/db"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	if !cfg.Server.Debug {
		util.StartSentry(cfg.SentryDSN, cfg.ReleaseFilePath)
		defer sentry.Flush(2 * time.Second)
	}
	db.MustConnect(cfg.Db)
	medsengerscalesbot.NewServer(cfg.Server).Listen()
}

