package main

import (
	"fmt"

	"github.com/tikhonp/medsenger-scales-bot/db"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

func main() {
	cfg := util.LoadConfigFromEnv()
	fmt.Print(
		db.DataSourceName(cfg.DB),
	)
}
