package bayes

import (
	"math"
	"strings"
	"unicode"

	"github.com/antfie/FoxBot/db"
)

const minTrainingExamples = 30

type Classifier struct {
	db *db.DB
}

func NewClassifier(db *db.DB) *Classifier {
	return &Classifier{db: db}
}

func (c *Classifier) Train(feedGroup, text string, relevant bool) {
	for _, word := range Tokenize(text) {
		c.db.BayesUpsertWord(feedGroup, word, relevant)
	}

	c.db.BayesIncrementStats(feedGroup, relevant)
}

func (c *Classifier) Untrain(feedGroup, text string, relevant bool) {
	for _, word := range Tokenize(text) {
		c.db.BayesDecrementWord(feedGroup, word, relevant)
	}

	c.db.BayesDecrementStats(feedGroup, relevant)
}

func (c *Classifier) Score(feedGroup, text string) float64 {
	words := Tokenize(text)

	if len(words) == 0 {
		return 0.5
	}

	relevant, irrelevant := c.db.BayesGetStats(feedGroup)
	total := relevant + irrelevant

	if total == 0 {
		return 0.5
	}

	wordCounts := c.db.BayesGetWordCounts(feedGroup)
	vocabSize := len(wordCounts)

	// Log prior probabilities
	logPriorRelevant := math.Log(float64(relevant) / float64(total))
	logPriorIrrelevant := math.Log(float64(irrelevant) / float64(total))

	logRelevant := logPriorRelevant
	logIrrelevant := logPriorIrrelevant

	for _, word := range words {
		counts, exists := wordCounts[word]

		var wordRelevant, wordIrrelevant int
		if exists {
			wordRelevant = counts[0]
			wordIrrelevant = counts[1]
		}

		// Laplace smoothing
		logRelevant += math.Log(float64(wordRelevant+1) / float64(relevant+vocabSize+1))
		logIrrelevant += math.Log(float64(wordIrrelevant+1) / float64(irrelevant+vocabSize+1))
	}

	// Convert from log space to probability using log-sum-exp for numerical stability
	maxLog := math.Max(logRelevant, logIrrelevant)
	logSum := maxLog + math.Log(math.Exp(logRelevant-maxLog)+math.Exp(logIrrelevant-maxLog))

	return math.Exp(logRelevant - logSum)
}

func (c *Classifier) IsReady(feedGroup string) bool {
	relevant, irrelevant := c.db.BayesGetStats(feedGroup)
	return relevant+irrelevant >= minTrainingExamples
}

func Tokenize(text string) []string {
	text = strings.ToLower(text)

	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	var result []string

	for _, word := range words {
		if len(word) >= 3 {
			result = append(result, word)
		}
	}

	return result
}
