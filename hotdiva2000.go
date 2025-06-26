// Package hotdiva2000 provides a human-readable random string generator.
package hotdiva2000

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:generate sort -u modifiers.txt -o modifiers.txt
//go:generate sort -u nouns.txt -o nouns.txt
//go:generate sort -u prefix.txt -o prefix.txt
//go:generate sort -u suffix.txt -o suffix.txt

const (
	prefixThreshold = 0.2
	suffixThreshold = 0.2
)

//go:embed prefix.txt
var prefixData string

//go:embed modifiers.txt
var modifierData string

//go:embed nouns.txt
var nounData string

//go:embed suffix.txt
var suffixData string

var (
	prefixes  []string = strings.Split(prefixData, "\n")
	modifiers []string = strings.Split(modifierData, "\n")
	nouns     []string = strings.Split(nounData, "\n")
	suffixes  []string = strings.Split(suffixData, "\n")

	// These start with vowels but should not be preceded with "an". Exceptions
	// will be checked as prefixes, so cases like "uptopia" will also over
	// "uptopian".
	anExceptions = []string{"unix", "utopia"}
)

func startsWithVowel(s string) bool {
	s = strings.ToLower(s)
	for _, e := range anExceptions {
		if strings.HasPrefix(s, e) {
			return false
		}
	}
	vowels := []string{"a", "e", "i", "o", "u"}
	for _, v := range vowels {
		if strings.HasPrefix(s, v) {
			return true
		}
	}
	return false
}

// Look for places where "a" should be "an" and correct accordingly.
func fixArticles(sentence string) string {
	words := strings.Split(sentence, " ")
	for i, word := range words {
		if strings.ToLower(word) == "a" && i < len(words)-1 {
			nextWord := words[i+1]
			if startsWithVowel(nextWord) {
				words[i] = "an"
			}
		}
	}

	return strings.Join(words, " ")
}

// Generate returns a random string.
func Generate() string {
	var (
		prefix = ""
		suffix = ""
	)

	if rand.Float64() < prefixThreshold {
		prefix = prefixes[rand.Intn(len(prefixes)-1)] + " "
	}
	if rand.Float64() < suffixThreshold {
		suffix = " " + suffixes[rand.Intn(len(suffixes)-1)]
	}

	mod := modifiers[rand.Intn(len(modifiers)-1)]
	noun := nouns[rand.Intn((len(nouns) - 1))]

	output := fixArticles(prefix + mod + " " + noun + suffix)
	return strings.ToLower(strings.ReplaceAll(output, " ", "-"))
}

// GenerateN returns a given number of random strings.
func GenerateN(n int) []string {
	r := make([]string, n)
	for i := 0; i < n; i++ {
		r[i] = Generate()
	}
	return r
}

// Possibilities return a number of possible strings produced.
func Possibilities() (int, int) {
	low := len(modifiers) * len(nouns)
	high := low * len(prefixes) * len(suffixes)
	return low, high
}
