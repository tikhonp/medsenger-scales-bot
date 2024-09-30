package main

import (
	"context"
	"fmt"

	"github.com/tikhonp/medsenger-scales-bot/config"
	"github.com/tikhonp/medsenger-scales-bot/db"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "pkl/local/config.pkl")
	if err != nil {
		panic(err)
	}
	fmt.Print(
		db.DataSourceName(cfg.Db),
	)
}

