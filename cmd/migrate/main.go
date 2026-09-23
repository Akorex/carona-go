package main

import (
	"fmt"
	"log"
	"os"

	"caronago/internal/platform/config"
	"caronago/internal/platform/db"

	"github.com/pressly/goose"
)

const migrationsDir = "migrations"

func main() {
	// verify that a command was passed
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go <command> [args]")
		fmt.Println("\nCommands:")
		fmt.Println("  up                   Run all pending migrations")
		fmt.Println("  down                 Roll back a single migration")
		fmt.Println("  status               Dump migration status")
		fmt.Println("  create <name> sql    Create a new migration file")
		fmt.Println("  reset                Roll back all migrations")
		os.Exit(1)
	}

	command := os.Args[1]
	var cmdArgs []string
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gormDB, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Run(command, sqlDB, migrationsDir, cmdArgs...); err != nil {
		log.Fatalf("goose run failed: %v", err)
	}

}
