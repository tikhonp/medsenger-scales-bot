package main

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	medsengerscalesbot "github.com/tikhonp/medsenger-scales-bot"
	"github.com/tikhonp/medsenger-scales-bot/db"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

func main() {
	cfg := util.LoadConfigFromEnv()
	if !cfg.Server.Debug {
		err := util.StartSentry(cfg.SentryDSN)
		if err != nil {
			log.Fatalln(err)
		}
		defer sentry.Flush(2 * time.Second)
	}
	db.MustConnect(cfg.DB)
	medsengerscalesbot.NewServer(cfg.Server).Listen()
}
