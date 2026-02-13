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
		log.Print(err)
		return nil
	}

	return rows
}

func (db *DB) Exec(query string, args ...any) {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.db.Exec(query, args...)

	if err != nil {
		log.Print(err)
	}
}

func (db *DB) insert(query string, args ...any) bool {
	result, err := db.db.Exec(query, args...)

	if err != nil {
		log.Print(err)
		return false
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Print(err)
		return false
	}

	return rowsAffected > 0
}

func (db *DB) IsRSSLinkInDB(link string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.db.Query("SELECT 1 FROM rss WHERE link = ? LIMIT 1", link)

	if err != nil {
		log.Print(err)
		return true // Assume exists to avoid duplicate insert attempts
	}

	found := rows.Next()
	err = rows.Err()

	if err != nil {
		log.Print(err)
		_ = rows.Close()
		return true
	}

	err = rows.Close()

	if err != nil {
		log.Print(err)
	}

	if !found {
		db.insert("INSERT INTO rss(link) VALUES (?)", link)
	}

	return found
}

func (db *DB) QueueTelegramNotification(message string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	success := db.insert("INSERT INTO telegram_notification(message) VALUES (?)", message)

	if !success {
		log.Print("Could not queue telegram notification")
	}
}

func (db *DB) ConsumeTelegramNotificationQueue() []string {
	db.mu.Lock()
	defer db.mu.Unlock()

	var results []string

	rows, err := db.db.Query("SELECT message FROM telegram_notification ORDER BY created")

	if err != nil {
		log.Print(err)
		return results
	}

	for rows.Next() {
		var value string
		err = rows.Scan(&value)

		if err != nil {
			log.Print(err)
			continue
		}

		if !utils.IsStringInArray(value, results) {
			results = append(results, value)
		}
	}

	err = rows.Err()

	if err != nil {
		log.Print(err)
	}

	err = rows.Close()

	if err != nil {
		log.Print(err)
	}

	_, err = db.db.Exec("DELETE FROM telegram_notification")

	if err != nil {
		log.Print(err)
	}

	return results
}

func (db *DB) QueueSlackNotification(message string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	success := db.insert("INSERT INTO slack_notification(message) VALUES (?)", message)

	if !success {
		log.Print("Could not queue slack notification")
	}
}

func (db *DB) ConsumeSlackNotificationQueue() []string {
	db.mu.Lock()
	defer db.mu.Unlock()

	var results []string

	rows, err := db.db.Query("SELECT message FROM slack_notification ORDER BY created")

	if err != nil {
		log.Print(err)
		return results
	}

	for rows.Next() {
		var value string
		err = rows.Scan(&value)

		if err != nil {
			log.Print(err)
			continue
		}

		if !utils.IsStringInArray(value, results) {
			results = append(results, value)
		}
	}

	err = rows.Err()

	if err != nil {
		log.Print(err)
	}

	err = rows.Close()

	if err != nil {
		log.Print(err)
	}

	_, err = db.db.Exec("DELETE FROM slack_notification")

	if err != nil {
		log.Print(err)
	}

	return results
}
