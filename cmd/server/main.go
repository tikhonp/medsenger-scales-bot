package main

import (
	"context"

	medsengerscalesbot "github.com/tikhonp/medsenger-scales-bot"
	"github.com/tikhonp/medsenger-scales-bot/config"
	"github.com/tikhonp/medsenger-scales-bot/db"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	if !cfg.Server.Debug {
		// TODO: Setup sentry
		// util.StartSentry(cfg.SentryDSN, cfg.ReleaseFilePath)
		// defer sentry.Flush(2 * time.Second)
	}
	db.MustConnect(cfg.Db)
	medsengerscalesbot.NewServer(cfg.Server).Listen()
}

