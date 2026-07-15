// Command manage is a small CLI for operational tasks.
//
//	manage -c print-db-string   # print the postgres DSN
//	manage -c migrate-up        # apply all pending migrations (run at container start)
//	manage -c migrate-down      # roll back a single migration
//	manage -c migrate-status    # print migration status
//	manage -c migrate-reset     # roll back all migrations
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tikhonp/medsenger-scales-bot/db"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

type command string

const (
	// PrintDBString prints the postgres connection string.
	PrintDBString command = "print-db-string"

	// MigrateUp applies all pending database migrations.
	MigrateUp command = "migrate-up"

	// MigrateDown rolls back a single database migration.
	MigrateDown command = "migrate-down"

	// MigrateStatus prints the database migration status.
	MigrateStatus command = "migrate-status"

	// MigrateReset rolls back all database migrations.
	MigrateReset command = "migrate-reset"
)

func (c *command) Set(value string) error {
	switch command(value) {
	case PrintDBString, MigrateUp, MigrateDown, MigrateStatus, MigrateReset:
		*c = command(value)
		return nil
	default:
		return fmt.Errorf("invalid command %s", value)
	}
}

func (c *command) String() string {
	return string(*c)
}

func main() {
	var cmd command
	const usage = "command to run. Available commands: print-db-string, migrate-up, migrate-down, migrate-status, migrate-reset"
	flag.Var(&cmd, "command", usage)
	flag.Var(&cmd, "c", usage+" (shorthand)")
	flag.Parse()

	cfg := util.LoadConfigFromEnv()

	switch cmd {
	case PrintDBString:
		fmt.Print(db.DataSourceName(cfg.DB))
	case MigrateUp:
		mustMigrate(cfg.DB, "up")
	case MigrateDown:
		mustMigrate(cfg.DB, "down")
	case MigrateStatus:
		mustMigrate(cfg.DB, "status")
	case MigrateReset:
		mustMigrate(cfg.DB, "reset")
	default:
		fmt.Println("Invalid arguments")
	}
}

func mustMigrate(cfg *util.Database, command string) {
	if err := db.Migrate(cfg, command); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
}
