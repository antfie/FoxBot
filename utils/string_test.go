package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringContainsWordIgnoreCase(t *testing.T) {
	keywords := []string{"apple", "banana", "cherry"}
	result := StringContainsWordIgnoreCase("oapple0.chErry.", keywords)
	assert.Equal(t, "cherry", result)
}
