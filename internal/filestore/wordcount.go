package filestore

import (
	"regexp"
	"sort"
	"strings"
)

var wordRegex *regexp.Regexp

func init() {
	wordRegex = regexp.MustCompile("[^a-zA-Z0-9]+")
}

func WordCount(content string) int {
	return len(strings.Fields(content))
}

func FrequentWords(allText string, limit int, order string) map[string]int {
	if allText == "" {
		return make(map[string]int)
	}

	processedString := wordRegex.ReplaceAllString(allText, " ")

	words := strings.Fields(processedString)
	wordFrequency := make(map[string]int)
	for _, word := range words {
		wordFrequency[strings.ToLower(word)]++
	}

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range wordFrequency {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		if order == "dsc" {
			return ss[i].Value > ss[j].Value
		}
		return ss[i].Value < ss[j].Value
	})

	if limit > len(ss) {
		limit = len(ss)
	}

	result := make(map[string]int)
	for _, kv := range ss[:limit] {
		result[kv.Key] = kv.Value
	}

	return result
}
