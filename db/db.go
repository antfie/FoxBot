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

	db.SetMaxOpenConns(1)

	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA foreign_keys = ON",
		"PRAGMA cache_size = -2000",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			log.Panicf("failed to execute %s: %v", p, err)
		}
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

// Bayes methods

func (db *DB) BayesUpsertWord(feedGroup, word string, relevant bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if relevant {
		db.insert("INSERT INTO bayes_model (feed_group, word, relevant, irrelevant) VALUES (?, ?, 1, 0) ON CONFLICT(feed_group, word) DO UPDATE SET relevant = relevant + 1", feedGroup, word)
	} else {
		db.insert("INSERT INTO bayes_model (feed_group, word, relevant, irrelevant) VALUES (?, ?, 0, 1) ON CONFLICT(feed_group, word) DO UPDATE SET irrelevant = irrelevant + 1", feedGroup, word)
	}
}

func (db *DB) BayesGetWordCounts(feedGroup string) map[string][2]int {
	db.mu.Lock()
	defer db.mu.Unlock()

	result := make(map[string][2]int)

	rows, err := db.db.Query("SELECT word, relevant, irrelevant FROM bayes_model WHERE feed_group = ?", feedGroup)

	if err != nil {
		log.Print(err)
		return result
	}

	for rows.Next() {
		var word string
		var relevant, irrelevant int
		err = rows.Scan(&word, &relevant, &irrelevant)

		if err != nil {
			log.Print(err)
			continue
		}

		result[word] = [2]int{relevant, irrelevant}
	}

	if err = rows.Err(); err != nil {
		log.Print(err)
	}

	if err = rows.Close(); err != nil {
		log.Print(err)
	}

	return result
}

func (db *DB) BayesIncrementStats(feedGroup string, relevant bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if relevant {
		db.insert("INSERT INTO bayes_stats (feed_group, relevant, irrelevant) VALUES (?, 1, 0) ON CONFLICT(feed_group) DO UPDATE SET relevant = relevant + 1", feedGroup)
	} else {
		db.insert("INSERT INTO bayes_stats (feed_group, relevant, irrelevant) VALUES (?, 0, 1) ON CONFLICT(feed_group) DO UPDATE SET irrelevant = irrelevant + 1", feedGroup)
	}
}

func (db *DB) BayesDecrementWord(feedGroup, word string, relevant bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if relevant {
		db.insert("UPDATE bayes_model SET relevant = MAX(relevant - 1, 0) WHERE feed_group = ? AND word = ?", feedGroup, word)
	} else {
		db.insert("UPDATE bayes_model SET irrelevant = MAX(irrelevant - 1, 0) WHERE feed_group = ? AND word = ?", feedGroup, word)
	}
}

func (db *DB) BayesDecrementStats(feedGroup string, relevant bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if relevant {
		db.insert("UPDATE bayes_stats SET relevant = MAX(relevant - 1, 0) WHERE feed_group = ?", feedGroup)
	} else {
		db.insert("UPDATE bayes_stats SET irrelevant = MAX(irrelevant - 1, 0) WHERE feed_group = ?", feedGroup)
	}
}

func (db *DB) BayesGetStats(feedGroup string) (int, int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.db.QueryRow("SELECT relevant, irrelevant FROM bayes_stats WHERE feed_group = ?", feedGroup)

	var relevant, irrelevant int
	err := row.Scan(&relevant, &irrelevant)

	if err != nil {
		return 0, 0
	}

	return relevant, irrelevant
}

func (db *DB) BayesSaveArticle(hash, feedGroup, title string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.insert("INSERT OR IGNORE INTO bayes_article (hash, feed_group, title) VALUES (?, ?, ?)", hash, feedGroup, title)
}

func (db *DB) BayesGetArticle(hash string) (feedGroup, title, label string, found bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.db.QueryRow("SELECT feed_group, title, COALESCE(label, '') FROM bayes_article WHERE hash = ?", hash)

	err := row.Scan(&feedGroup, &title, &label)

	if err != nil {
		return "", "", "", false
	}

	return feedGroup, title, label, true
}

func (db *DB) BayesSetArticleLabel(hash, label string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.insert("UPDATE bayes_article SET label = ? WHERE hash = ?", label, hash)
}

func (db *DB) BayesCleanupOldArticles() {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.db.Exec("DELETE FROM bayes_article WHERE created < date('now', '-30 day')")

	if err != nil {
		log.Print(err)
	}
}

func (db *DB) GetTelegramState(key string) string {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.db.QueryRow("SELECT value FROM telegram_state WHERE key = ?", key)

	var value string
	err := row.Scan(&value)

	if err != nil {
		return ""
	}

	return value
}

func (db *DB) SetTelegramState(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.insert("INSERT INTO telegram_state (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?", key, value, value)
}

// HTTP cache methods

func (db *DB) GetHTTPCache(url string) (etag, lastModified string, failCount int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.db.QueryRow("SELECT etag, last_modified, fail_count FROM http_cache WHERE url = ?", url)

	err := row.Scan(&etag, &lastModified, &failCount)

	if err != nil {
		return "", "", 0
	}

	return etag, lastModified, failCount
}

func (db *DB) SetHTTPCache(url, etag, lastModified string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.insert("INSERT INTO http_cache (url, etag, last_modified, fail_count) VALUES (?, ?, ?, 0) ON CONFLICT(url) DO UPDATE SET etag = ?, last_modified = ?, fail_count = 0", url, etag, lastModified, etag, lastModified)
}

func (db *DB) IncrementHTTPCacheFailCount(url string) int {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.insert("INSERT INTO http_cache (url, fail_count) VALUES (?, 1) ON CONFLICT(url) DO UPDATE SET fail_count = fail_count + 1", url)

	row := db.db.QueryRow("SELECT fail_count FROM http_cache WHERE url = ?", url)

	var failCount int
	err := row.Scan(&failCount)

	if err != nil {
		return 0
	}

	return failCount
}

// Weather methods

func (db *DB) HasWeatherBeenNotifiedToday(location string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.db.QueryRow("SELECT 1 FROM weather_notification WHERE location = ? AND last_notified = date('now')", location)

	var dummy int
	err := row.Scan(&dummy)

	return err == nil
}

func (db *DB) SetWeatherNotified(location string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.insert("INSERT INTO weather_notification (location, last_notified) VALUES (?, date('now')) ON CONFLICT(location) DO UPDATE SET last_notified = date('now')", location)
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
