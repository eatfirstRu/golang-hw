package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var taskWithAsteriskIsCompleted = true

func Top10(inStr string) []string {
	var re regexp.Regexp
	var cntWord = make(map[string]int)
	var resultWords = make([]string, 0, 10)

	if taskWithAsteriskIsCompleted {
		inStr = strings.ToLower(inStr)
		re = *regexp.MustCompile("[\"!?,.:;]|\n\t|\t| -|- ")
	} else {
		re = *regexp.MustCompile("[\n\t]")
	}
	inStr = re.ReplaceAllString(inStr, " ")

	words := strings.Split(inStr, " ")

	for _, word := range words {
		if word != "" {
			cntWord[word]++
		}
	}

	words = nil

	for word := range cntWord {
		nmb := fmt.Sprintf("%06d", cntWord[word]-1000000)
		words = append(words, nmb+"<<\\t>>"+word)
	}
	sort.Strings(words)

	lenWords := len(words)
	if lenWords > 10 {
		lenWords = 10
	}

	for i := 0; i < lenWords; i++ {
		_, wrd, _ := strings.Cut(words[i], "<<\\t>>")
		resultWords = append(resultWords, wrd)
	}
	return resultWords
}
