package hotdiva2000

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	result := Generate()
	if result == "" {
		t.Fatal("Generate returned empty string")
	}
	if strings.Contains(result, " ") {
		t.Errorf("Generate result should not contain spaces, got: %q", result)
	}
}

func TestGenerateN(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"single", 1, 1},
		{"five", 5, 5},
		{"ten", 10, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := GenerateN(tt.n)
			if len(results) != tt.want {
				t.Errorf("GenerateN(%d) returned %d results, want %d", tt.n, len(results), tt.want)
			}
			for i, r := range results {
				if r == "" {
					t.Errorf("GenerateN(%d)[%d] is empty", tt.n, i)
				}
			}
		})
	}
}

func TestGenerateNZero(t *testing.T) {
	// n < 1 should clamp to 1
	results := GenerateN(0)
	if len(results) != 1 {
		t.Errorf("GenerateN(0) returned %d results, want 1", len(results))
	}
}

func TestGenerateWithOptions(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{
			name: "no prefix or suffix",
			opts: Options{PrefixThreshold: 0, SuffixThreshold: 0, Results: 10},
		},
		{
			name: "always prefix and suffix",
			opts: Options{PrefixThreshold: 1, SuffixThreshold: 1, Results: 10},
		},
		{
			name: "default thresholds",
			opts: Options{PrefixThreshold: 0.2, SuffixThreshold: 0.2, Results: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := GenerateWithOptions(tt.opts)
			if len(results) != tt.opts.Results {
				t.Errorf("got %d results, want %d", len(results), tt.opts.Results)
			}
			for i, r := range results {
				if r == "" {
					t.Errorf("result[%d] is empty", i)
				}
				if strings.Contains(r, " ") {
					t.Errorf("result[%d] contains spaces: %q", i, r)
				}
			}
		})
	}
}

func TestGenerateWithOptionsAlwaysPrefixSuffix(t *testing.T) {
	// With threshold 1.0, every result should have at least 3 hyphens
	// (prefix-modifier-noun-suffix)
	results := GenerateWithOptions(Options{
		PrefixThreshold: 1,
		SuffixThreshold: 1,
		Results:         50,
	})
	for i, r := range results {
		parts := strings.Split(r, "-")
		if len(parts) < 4 {
			t.Errorf("result[%d] = %q has %d parts, expected at least 4 with prefix+suffix", i, r, len(parts))
		}
	}
}

func TestGenerateWithOptionsNoPrefixSuffix(t *testing.T) {
	// With threshold 0, results should have exactly 2+ parts (modifier-noun),
	// some words may themselves contain hyphens so we check for at least 2
	results := GenerateWithOptions(Options{
		PrefixThreshold: 0,
		SuffixThreshold: 0,
		Results:         50,
	})
	for i, r := range results {
		parts := strings.Split(r, "-")
		if len(parts) < 2 {
			t.Errorf("result[%d] = %q has fewer than 2 parts", i, r)
		}
	}
}

func TestPossibilities(t *testing.T) {
	low, high := Possibilities()
	if low <= 0 {
		t.Errorf("low possibilities should be > 0, got %d", low)
	}
	if high <= 0 {
		t.Errorf("high possibilities should be > 0, got %d", high)
	}
	if high < low {
		t.Errorf("high (%d) should be >= low (%d)", high, low)
	}
}

func TestPossibilitiesWithOptions(t *testing.T) {
	tests := []struct {
		name string
		opts Options
	}{
		{"no extras", Options{PrefixThreshold: 0, SuffixThreshold: 0}},
		{"always extras", Options{PrefixThreshold: 1, SuffixThreshold: 1}},
		{"partial", Options{PrefixThreshold: 0.5, SuffixThreshold: 0.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			low, high := PossibilitiesWithOptions(tt.opts)
			if low <= 0 {
				t.Errorf("low should be > 0, got %d", low)
			}
			if high < low {
				t.Errorf("high (%d) should be >= low (%d)", high, low)
			}
		})
	}
}

func TestPossibilitiesWithOptionsMonotonicity(t *testing.T) {
	// More prefix/suffix threshold should increase high possibilities
	_, highNone := PossibilitiesWithOptions(Options{PrefixThreshold: 0, SuffixThreshold: 0})
	_, highFull := PossibilitiesWithOptions(Options{PrefixThreshold: 1, SuffixThreshold: 1})
	if highFull < highNone {
		t.Errorf("full thresholds high (%d) should be >= no thresholds high (%d)", highFull, highNone)
	}
}

func TestFixArticles(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"a apple", "an apple"},
		{"a banana", "a banana"},
		{"a elephant", "an elephant"},
		{"a orange", "an orange"},
		{"a umbrella", "an umbrella"},
		{"a ice cream", "an ice cream"},
		{"the apple", "the apple"},
		{"a unix guru", "a unix guru"},       // exception
		{"a utopian place", "a utopian place"}, // exception
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := fixArticles(tt.input)
			if got != tt.want {
				t.Errorf("fixArticles(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStartsWithVowel(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"apple", true},
		{"banana", false},
		{"elephant", true},
		{"ice", true},
		{"orange", true},
		{"umbrella", true},
		{"unix", false},    // exception
		{"utopia", false},  // exception
		{"utopian", false}, // exception prefix match
		{"Elephant", true}, // case insensitive
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := startsWithVowel(tt.input)
			if got != tt.want {
				t.Errorf("startsWithVowel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateUniqueness(t *testing.T) {
	// Generate a batch and check that not all results are identical
	results := GenerateN(100)
	seen := make(map[string]bool)
	for _, r := range results {
		seen[r] = true
	}
	if len(seen) < 2 {
		t.Error("100 generated strings are all identical, expected variety")
	}
}

func TestGenerateLowercase(t *testing.T) {
	// All results should be lowercase
	results := GenerateN(100)
	for i, r := range results {
		if r != strings.ToLower(r) {
			t.Errorf("result[%d] = %q is not lowercase", i, r)
		}
	}
}

func TestGenerateNoSpaces(t *testing.T) {
	results := GenerateN(100)
	for i, r := range results {
		if strings.Contains(r, " ") {
			t.Errorf("result[%d] = %q contains spaces", i, r)
		}
	}
}

func TestWordListsLoaded(t *testing.T) {
	if len(prefixes) == 0 {
		t.Error("prefixes word list is empty")
	}
	if len(modifiers) == 0 {
		t.Error("modifiers word list is empty")
	}
	if len(nouns) == 0 {
		t.Error("nouns word list is empty")
	}
	if len(suffixes) == 0 {
		t.Error("suffixes word list is empty")
	}
}

func TestThresholdClamping(t *testing.T) {
	// Thresholds > 1 or < 0 should be clamped and not panic
	results := GenerateWithOptions(Options{
		PrefixThreshold: 2.0,
		SuffixThreshold: -1.0,
		Results:         5,
	})
	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
}

func BenchmarkGenerate(b *testing.B) {
	for b.Loop() {
		Generate()
	}
}

func BenchmarkGenerateN100(b *testing.B) {
	for b.Loop() {
		GenerateN(100)
	}
}
