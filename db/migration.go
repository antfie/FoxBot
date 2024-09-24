package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func runMigrations(db *sql.DB) {
	currentMigration := getDbMigrationVersion(db)

	for {
		currentMigration = currentMigration + 1

		migrationFile := path.Join("migrations", fmt.Sprintf("%03d.sql", currentMigration))

		data, err := migrationFiles.ReadFile(migrationFile)

		if err != nil {
			// This file does not exist, so no more migrations to run
			return
		}

		log.Printf("Performing DB migration %d...", currentMigration)

		tx, err := db.Begin()

		if err != nil {
			log.Fatal(err)
		}

		// Is this the first migration?
		if currentMigration == 1 {
			_, err = tx.Query("CREATE TABLE migration (version INTEGER PRIMARY KEY)")

			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = tx.Query(string(data))

		if err != nil {
			log.Fatal(err)
		}

		_, err = tx.Query("INSERT INTO migration (version) VALUES (?)", currentMigration)

		if err != nil {
			log.Fatal(err)
		}

		err = tx.Commit()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func getDbMigrationVersion(db *sql.DB) int {
	rows, err := db.Query("SELECT 1 FROM sqlite_schema WHERE type='table' AND name='migration' LIMIT 1")

	if err != nil {
		log.Panic(err)
	}

	found := rows.Next()

	err = rows.Close()

	if err != nil {
		log.Panic(err)
	}

	if !found {
		return 0
	}

	rows, err = db.Query("SELECT version FROM migration ORDER BY version DESC LIMIT 1")

	if err != nil {
		log.Panic(err)
	}

	found = rows.Next()

	if !found {
		return 0
	}

	currentMigration := 0
	err = rows.Scan(&currentMigration)

	if err != nil {
		log.Panic(err)
	}

	err = rows.Close()

	if err != nil {
		log.Panic(err)
	}

	return currentMigration
}
