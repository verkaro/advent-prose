package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/verkaro/advent-prose"
)

func main() {
	// Define flags to control the ventilation configuration.
	inputFile := flag.String("i", "", "Input file path. If not specified, reads from stdin.")
	sentenceBreak := flag.Bool("sentence-break", true, "Break lines at the end of sentences.")
	maxLineLength := flag.Int("max-line-length", 0, "Soft limit for line length (only used if sentence-break is false).")
	pSpacing := flag.String("paragraph-spacing", "single", "Paragraph spacing ('single' or 'blank-line').")
	respectMaxLine := flag.Bool("respect-max-line-length", false, "Respect max line length for soft wrapping.")
	abbrevs := flag.String("abbreviations", "", "Comma-separated list of custom abbreviations (e.g., \"No.,Fig.\").")

	// Custom usage message to provide more context and examples.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Advent Prose Ventilator: A tool to reflow Markdown prose.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Processes Markdown from standard input or an input file and prints to standard output.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Process from stdin\n")
		fmt.Fprintf(os.Stderr, "  cat my_document.md | %s\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  # Process from a file\n")
		fmt.Fprintf(os.Stderr, "  %s -i my_document.md\n", os.Args[0])
	}

	flag.Parse()

	var inputBytes []byte
	var err error

	// Read from the specified input file or fall back to stdin.
	if *inputFile != "" {
		inputBytes, err = os.ReadFile(*inputFile)
		if err != nil {
			log.Fatalf("Error reading from file %q: %v", *inputFile, err)
		}
	} else {
		inputBytes, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
	}

	// Build the configuration from the flags.
	cfg := advent.Config{
		SentenceBreak:        *sentenceBreak,
		MaxLineLength:        *maxLineLength,
		ParagraphSpacing:     *pSpacing,
		RespectMaxLineLength: *respectMaxLine,
	}

	// If custom abbreviations are provided, parse them into the config map.
	if *abbrevs != "" {
		cfg.Abbreviations = make(map[string]bool)
		parts := strings.Split(*abbrevs, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				cfg.Abbreviations[trimmed] = true
			}
		}
	}

	// Ventilate the input using the library.
	output, err := advent.Ventilate(string(inputBytes), cfg)
	if err != nil {
		log.Fatalf("Error ventilating text: %v", err)
	}

	// Write the result to standard output.
	fmt.Println(output)
}

