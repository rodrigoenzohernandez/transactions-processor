package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/repository"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

var log = logger.GetLogger("db-scripts")

func runMigrations(migrationStr string, operation string) {

	cwd, _ := os.Getwd()

	migrationsPath := filepath.Join(cwd, "db/migrations")

	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.New(sourceURL, migrationStr)

	if err != nil {
		log.Error(fmt.Sprintf("Error creating migrate instance: %v", err))

	}

	switch operation {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Error(fmt.Sprintf("Invalid migration operation: %s", err))
		return
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Error(fmt.Sprintf("Error running migrations: %v", err))

	}

	log.Info("Migrations applied successfully!")

}

func main() {

	if len(os.Args) < 2 {
		log.Error("Insert the operation <up|down>")
		return
	}

	operation := os.Args[1]

	if operation != "up" && operation != "down" {
		log.Error("Invalid operation")

		return
	}

	db, migrationStr := repository.Connect()

	defer repository.Disconnect(db)

	runMigrations(migrationStr, operation)
}
