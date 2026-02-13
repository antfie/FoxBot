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

func TestStringContainsWordIgnoreCaseNoMatch(t *testing.T) {
	keywords := []string{"apple", "banana"}
	result := StringContainsWordIgnoreCase("there is no fruit here", keywords)
	assert.Equal(t, "", result)
}

func TestStringContainsWordIgnoreCaseEmptyArray(t *testing.T) {
	result := StringContainsWordIgnoreCase("anything", []string{})
	assert.Equal(t, "", result)
}

func TestIsStringInArray(t *testing.T) {
	arr := []string{"a", "b", "c"}
	assert.True(t, IsStringInArray("b", arr))
	assert.False(t, IsStringInArray("d", arr))
	assert.False(t, IsStringInArray("B", arr))
}

func TestIsStringInArrayEmpty(t *testing.T) {
	assert.False(t, IsStringInArray("a", []string{}))
	assert.False(t, IsStringInArray("a", nil))
}

func TestMergeStringArrays(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"b", "c", "d"}
	result := MergeStringArrays(a, b)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result)
}

func TestMergeStringArraysNoDuplicates(t *testing.T) {
	a := []string{"x", "y"}
	b := []string{"x", "y"}
	result := MergeStringArrays(a, b)
	assert.Equal(t, []string{"x", "y"}, result)
}

func TestMergeStringArraysEmpty(t *testing.T) {
	result := MergeStringArrays([]string{}, []string{"a"})
	assert.Equal(t, []string{"a"}, result)

	result = MergeStringArrays([]string{"a"}, []string{})
	assert.Equal(t, []string{"a"}, result)
}

func TestPluralize(t *testing.T) {
	assert.Equal(t, "1 task", Pluralize("task", 1))
	assert.Equal(t, "0 tasks", Pluralize("task", 0))
	assert.Equal(t, "5 tasks", Pluralize("task", 5))
}

func TestStringToFloat(t *testing.T) {
	val, err := StringToFloat("3.14")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, val)

	val, err = StringToFloat("0.0")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, val)

	_, err = StringToFloat("notanumber")
	assert.Error(t, err)
}
