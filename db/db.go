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
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
}

// MustConnect creates a new in-memory SQLite database and initializes it with the schema.
func MustConnect(cfg *util.Database) {
	db = sqlx.MustConnect("postgres", DataSourceName(cfg))
}
