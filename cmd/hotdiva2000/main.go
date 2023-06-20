package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/hotdiva2000"
	"github.com/charmbracelet/x/exp/ordered"
	"github.com/dustin/go-humanize"
	"github.com/mattn/go-runewidth"
	flag "github.com/spf13/pflag"
)

func formatPossibilities() string {
	low, high := hotdiva2000.Possibilities()
	return fmt.Sprintf(
		"Minimum combinations: %s\nMaximum combinations: %s",
		humanize.Comma(int64(low)),
		humanize.Comma(int64(high)),
	)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS]", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\n\n%s\n\n", formatPossibilities())
	fmt.Fprintln(os.Stderr, "Options:")

	var flags []string
	var descs []string
	flag.VisitAll(func(f *flag.Flag) {
		if f.Shorthand == "" {
			flags = append(flags, fmt.Sprintf("--%s", f.Name))
		} else {
			flags = append(flags, fmt.Sprintf("-%s --%s", f.Shorthand, f.Name))
		}
		descs = append(descs, f.Usage)
	})

	var widestFlag int
	for _, f := range flags {
		w := runewidth.StringWidth(f)
		if w > widestFlag {
			widestFlag = w
		}
	}

	const gap = 2
	for i := range flags {
		fmt.Fprintf(os.Stderr, "  %-*s %s\n", widestFlag+gap, flags[i], descs[i])
	}
}

func main() {
	const (
		minResults     = 1
		maxResults     = 1000
		defaultResults = 1
	)

	var (
		showHelp bool
		results  int
	)

	flag.BoolVarP(&showHelp, "help", "h", false, "Show this help and exit")
	flag.IntVarP(&results, "results", "r", defaultResults, "Number of results to generate (deafult 1)")

	flag.CommandLine.SortFlags = false
	flag.Usage = usage
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(1)
	}

	results = ordered.Clamp(results, minResults, maxResults)

	for i := 0; i < results; i++ {
		fmt.Println(hotdiva2000.Generate())
	}
}
