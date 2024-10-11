package db

import (
	"database/sql"
	"github.com/antfie/FoxBot/utils"
	"log"
	"sync"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
	mu sync.Mutex
}

func NewDB(dbPath string) *DB {
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		log.Panic(err)
	}

	runMigrations(db)

	return &DB{db: db}
}

func (db *DB) Query(query string, args ...any) *sql.Rows {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.db.Query(query, args...)

	if err != nil {
		log.Panic(err)
	}

	return rows
}

func (db *DB) insert(query string, args ...any) bool {
	result, err := db.db.Exec(query, args...)

	if err != nil {
		log.Panic(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Panic(err)
	}

	return rowsAffected > 0
}

func (db *DB) IsRSSLinkInDB(link string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.db.Query("SELECT 1 FROM rss WHERE link = ? LIMIT 1", link)

	if err != nil {
		log.Panic(err)
	}

	found := rows.Next()
	err = rows.Err()

	if err != nil {
		log.Panic(err)
	}

	err = rows.Close()

	if err != nil {
		log.Panic(err)
	}

	if !found {
		db.insert("INSERT INTO rss(link) VALUES (?)", link)
	}

	return found
}

func (db *DB) QueueSlackNotification(message string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	success := db.insert("INSERT INTO slack_notification(message) VALUES (?)", message)

	if !success {
		log.Panic("Could not queue slack notification")
	}
}

func (db *DB) ConsumeSlackNotificationQueue() []string {
	db.mu.Lock()
	defer db.mu.Unlock()

	var results []string

	rows, err := db.db.Query("SELECT message FROM slack_notification ORDER BY created")

	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var value string
		err = rows.Scan(&value)

		if err != nil {
			log.Panic(err)
		}

		if !utils.IsStringInArray(value, results) {
			results = append(results, value)
		}
	}

	err = rows.Err()

	if err != nil {
		log.Panic(err)
	}

	err = rows.Close()

	if err != nil {
		log.Panic(err)
	}

	_, err = db.db.Exec("DELETE FROM slack_notification")

	if err != nil {
		log.Panic(err)
	}

	return results
}
