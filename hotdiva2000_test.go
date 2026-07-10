package hotdiva2000

import (
	"slices"
	"strings"
	"testing"
)

func containsUpper(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func TestFormatting(t *testing.T) {
	const n = 1000

	base := Options{
		PrefixThreshold: 1,
		SuffixThreshold: 1,
		Results:         n,
	}

	// FormatSlug produces lowercase, hyphenated output.
	slug := base
	slug.Formatting = FormatSlug
	for _, r := range GenerateWithOptions(slug) {
		if strings.Contains(r, " ") || containsUpper(r) {
			t.Errorf("FormatSlug: expected no spaces or uppercase, got %q", r)
		}
	}

	// FormatLowers produces lowercase output with spaces.
	lowers := base
	lowers.Formatting = FormatLowers
	for _, r := range GenerateWithOptions(lowers) {
		if r != strings.ToLower(r) {
			t.Errorf("FormatLowers: expected all lowercase, got %q", r)
		}
	}

	// FormatTitle preserves the original title case.
	title := base
	title.Formatting = FormatTitle
	if !slices.ContainsFunc(GenerateWithOptions(title), containsUpper) {
		t.Error("FormatTitle: expected title case to be preserved, but no uppercase found")
	}
}

func TestFormattingDefault(t *testing.T) {
	// A zero-value Formatting behaves like FormatSlug.
	opts := Options{
		PrefixThreshold: 1,
		SuffixThreshold: 1,
		Results:         100,
	}
	for _, r := range GenerateWithOptions(opts) {
		if strings.Contains(r, " ") || containsUpper(r) {
			t.Errorf("default Formatting: expected slug behavior, got %q", r)
		}
	}
}
