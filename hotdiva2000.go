// Package hotdiva2000 provides a human-readable random string generator.
package hotdiva2000

import (
	_ "embed"
	"math/rand"
	"strings"

	"github.com/charmbracelet/x/exp/ordered"
)

//go:generate sort -u modifiers.txt -o modifiers.txt
//go:generate sort -u nouns.txt -o nouns.txt
//go:generate sort -u prefix.txt -o prefix.txt
//go:generate sort -u suffix.txt -o suffix.txt

const (
	defaultPrefixThreshold = 0.2
	defaultSuffixThreshold = 0.2
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

// generate returns a random strings.
func generate(opts Options) []string {
	if opts.Results < 1 {
		opts.Results = 1
	}
	opts.PrefixThreshold = ordered.Clamp(opts.PrefixThreshold, 0, 1)
	opts.SuffixThreshold = ordered.Clamp(opts.SuffixThreshold, 0, 1)

	r := make([]string, opts.Results)

	for i := range r {
		var (
			prefix = ""
			suffix = ""
		)

		if opts.PrefixThreshold > 0 {
			if rand.Float64() < defaultPrefixThreshold {
				prefix = prefixes[rand.Intn(len(prefixes)-1)] + " "
			}
		}
		if opts.SuffixThreshold > 0 {
			if rand.Float64() < defaultSuffixThreshold {
				suffix = " " + suffixes[rand.Intn(len(suffixes)-1)]
			}
		}

		mod := modifiers[rand.Intn(len(modifiers)-1)]
		noun := nouns[rand.Intn((len(nouns) - 1))]

		output := fixArticles(prefix + mod + " " + noun + suffix)
		r[i] = strings.ToLower(strings.ReplaceAll(output, " ", "-"))
	}

	return r
}

// Options are options to customize output.
type Options struct {
	// Whether to show occasional additional prefix and suffix content. This
	// increases possibilities but can make strings longer.
	PrefixThreshold float64
	SuffixThreshold float64

	// Number of results to generate.
	Results int
}

// Generate returns a random string.
func Generate() string {
	return generate(Options{
		PrefixThreshold: defaultPrefixThreshold,
		SuffixThreshold: defaultSuffixThreshold,
	})[0]
}

// GenerateN returns a given number of random strings.
func GenerateN(n int) []string {
	return generate(Options{
		PrefixThreshold: defaultPrefixThreshold,
		SuffixThreshold: defaultSuffixThreshold,
		Results:         n,
	})
}

// Possibilities returns the number of possible strings produced.
func Possibilities() (int, int) {
	low := len(modifiers) * len(nouns)
	high := low * len(prefixes) * len(suffixes)
	return low, high
}

// GenerateWithOptions generates results against the given options.
func GenerateWithOptions(o Options) []string {
	return generate(o)
}

// PossibilitiesWithOptions returns the number of possible strings produced
// against the given options.
func PossibilitiesWithOptions(o Options) (int, int) {
	low := len(modifiers) * len(nouns)
	high := low

	o.PrefixThreshold = ordered.Clamp(o.PrefixThreshold, 0, 1)
	o.SuffixThreshold = ordered.Clamp(o.SuffixThreshold, 0, 1)

	if o.PrefixThreshold >= 1 {
		low *= len(prefixes)
	}
	if o.SuffixThreshold >= 1 {
		low *= len(suffixes)
	}

	if o.PrefixThreshold > 0 {
		high *= len(prefixes)
	}
	if o.SuffixThreshold > 0 {
		high *= len(suffixes)
	}

	return low, high
}
