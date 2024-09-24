package db

import (
	"database/sql"
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
	//db.SetMaxOpenConns(1)
	runMigrations(db)

	if err != nil {
		log.Panic(err)
	}

	return &DB{db: db}
}

func (db *DB) Query(query string, args ...any) *sql.Rows {
	rows, err := db.db.Query(query, args...)

	if err != nil {
		log.Panic(err)
	}

	return rows
}

func (db *DB) Insert(query string, args ...any) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

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
		db.Insert("INSERT INTO rss(link) VALUES (?)", link)
	}

	return found
}
