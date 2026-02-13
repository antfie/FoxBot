package bayes

import (
	"testing"

	"github.com/antfie/FoxBot/db"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *db.DB {
	t.Helper()
	return db.NewDB(":memory:")
}

func TestTokenize(t *testing.T) {
	assert.Equal(t, []string{"hello", "world"}, Tokenize("Hello World"))
	assert.Equal(t, []string{"npm", "malware", "found", "new", "package"}, Tokenize("npm malware found in new package"))
	assert.Equal(t, []string{"test123", "value"}, Tokenize("test123 value"))
}

func TestTokenizeDropsShortWords(t *testing.T) {
	result := Tokenize("I am a go dev")
	assert.Equal(t, []string{"dev"}, result)
}

func TestTokenizeEmpty(t *testing.T) {
	assert.Empty(t, Tokenize(""))
	assert.Empty(t, Tokenize("  "))
	assert.Empty(t, Tokenize("a b c"))
}

func TestTokenizePunctuation(t *testing.T) {
	result := Tokenize("hello, world! this is a test.")
	assert.Equal(t, []string{"hello", "world", "this", "test"}, result)
}

func TestClassifierIsReadyFalseWhenNew(t *testing.T) {
	d := setupTestDB(t)
	c := NewClassifier(d)
	assert.False(t, c.IsReady("test"))
}

func TestClassifierIsReadyAfterTraining(t *testing.T) {
	d := setupTestDB(t)
	c := NewClassifier(d)

	for i := 0; i < 20; i++ {
		c.Train("test", "relevant security article about npm malware", true)
	}
	for i := 0; i < 10; i++ {
		c.Train("test", "irrelevant sports football match results", false)
	}

	assert.True(t, c.IsReady("test"))
}

func TestClassifierScoreWithNoData(t *testing.T) {
	d := setupTestDB(t)
	c := NewClassifier(d)
	score := c.Score("test", "hello world")
	assert.Equal(t, 0.5, score)
}

func TestClassifierScoreEmptyText(t *testing.T) {
	d := setupTestDB(t)
	c := NewClassifier(d)
	c.Train("test", "some training data here", true)
	score := c.Score("test", "")
	assert.Equal(t, 0.5, score)
}

func TestClassifierTrainAndScore(t *testing.T) {
	d := setupTestDB(t)
	c := NewClassifier(d)

	// Train with security-related content as relevant
	for i := 0; i < 15; i++ {
		c.Train("security", "critical npm malware package detected supply chain attack", true)
		c.Train("security", "football match results weather forecast celebrity gossip", false)
	}

	// Score a security article - should be relevant
	securityScore := c.Score("security", "new npm malware package found in supply chain")
	assert.Greater(t, securityScore, 0.5)

	// Score a non-security article - should be irrelevant
	sportsScore := c.Score("security", "football match results and celebrity news today")
	assert.Less(t, sportsScore, 0.5)
}

func TestClassifierSeparateFeedGroups(t *testing.T) {
	d := setupTestDB(t)
	c := NewClassifier(d)

	for i := 0; i < 15; i++ {
		c.Train("security", "malware detected in npm package", true)
		c.Train("security", "weather forecast sunny today", false)
	}

	// Security group should be ready
	assert.True(t, c.IsReady("security"))

	// BBC group should not be ready
	assert.False(t, c.IsReady("bbc"))

	// Scoring against untrained group returns 0.5
	assert.Equal(t, 0.5, c.Score("bbc", "malware detected"))
}

func TestClassifierUnseenWords(t *testing.T) {
	d := setupTestDB(t)
	c := NewClassifier(d)

	for i := 0; i < 15; i++ {
		c.Train("test", "security vulnerability found", true)
		c.Train("test", "sports results today", false)
	}

	// Score with completely unseen words should still return a value
	score := c.Score("test", "completely novel unprecedented terminology")
	assert.Greater(t, score, 0.0)
	assert.Less(t, score, 1.0)
}
