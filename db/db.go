package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

// db is a global database.
//
// Yes, im dumb and i use global varibles for db.
// It's my second project in go, i think you can forgive me.
var db *sqlx.DB

func DataSourceName(cfg *util.Database) string {
	return fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", cfg.User, cfg.Dbname, cfg.Password, cfg.Host)
}

// MustConnect creates a new in-memory SQLite database and initializes it with the schema.
func MustConnect(cfg *util.Database) {
	db = sqlx.MustConnect("postgres", DataSourceName(cfg))
}
