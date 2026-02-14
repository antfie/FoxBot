package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "foxbot-test-*.db")
	assert.NoError(t, err)
	tmpFile.Close()

	database := NewDB(tmpFile.Name())

	return database, func() {
		os.Remove(tmpFile.Name())
	}
}

func TestIsRSSLinkInDB(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// First call should return false (not in DB) and insert
	assert.False(t, db.IsRSSLinkInDB("https://example.com/article1"))

	// Second call should return true (now in DB)
	assert.True(t, db.IsRSSLinkInDB("https://example.com/article1"))

	// Different link should return false
	assert.False(t, db.IsRSSLinkInDB("https://example.com/article2"))
}

func TestQueueAndConsumeSlackNotifications(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Empty queue
	messages := db.ConsumeSlackNotificationQueue()
	assert.Empty(t, messages)

	// Queue messages
	db.QueueSlackNotification("hello")
	db.QueueSlackNotification("world")

	// Consume returns all messages
	messages = db.ConsumeSlackNotificationQueue()
	assert.Equal(t, []string{"hello", "world"}, messages)

	// Queue is now empty
	messages = db.ConsumeSlackNotificationQueue()
	assert.Empty(t, messages)
}

func TestQueueSlackNotificationDeduplicates(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	db.QueueSlackNotification("duplicate")
	db.QueueSlackNotification("duplicate")
	db.QueueSlackNotification("unique")

	messages := db.ConsumeSlackNotificationQueue()
	assert.Equal(t, []string{"duplicate", "unique"}, messages)
}

func TestWeatherNotification(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Not notified yet
	assert.False(t, db.HasWeatherBeenNotifiedToday("Manchester"))

	// Mark as notified
	db.SetWeatherNotified("Manchester")

	// Now it should be true
	assert.True(t, db.HasWeatherBeenNotifiedToday("Manchester"))

	// Different location should still be false
	assert.False(t, db.HasWeatherBeenNotifiedToday("London"))
}

func TestExec(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Insert via IsRSSLinkInDB
	db.IsRSSLinkInDB("https://example.com/old")

	// Exec a delete
	db.Exec("DELETE FROM rss WHERE link = ?", "https://example.com/old")

	// Should be gone - next call inserts again and returns false
	assert.False(t, db.IsRSSLinkInDB("https://example.com/old"))
}
