package utils

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func StringContainsWordIgnoreCase(value string, array []string) string {
	pattern := fmt.Sprintf(`(?i)\b(?:%s)\b`, strings.Join(array, "|"))

	re, err := regexp.Compile(strings.ToLower(pattern))
	if err != nil {
		log.Panic(err)
	}

	matches := re.FindAllString(strings.ToLower(value), 1)

	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}

func MergeStringArrays(a, b []string) []string {
	output := a

	for _, x := range b {
		if !slices.Contains(output, x) {
			output = append(output, x)
		}
	}

	return output
}

func Pluralize(s string, count int) string {
	if count == 1 {
		return fmt.Sprintf("1 %s", s)
	}

	return fmt.Sprintf("%d %ss", count, s)
}

func StringToFloat(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}
