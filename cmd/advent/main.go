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
	sentenceBreak := flag.Bool("sentence-break", true, "Break lines at the end of sentences.")
	maxLineLength := flag.Int("max-line-length", 0, "Soft limit for line length (only used if sentence-break is false).")
	pSpacing := flag.String("paragraph-spacing", "single", "Paragraph spacing ('single' or 'blank-line').")
	respectMaxLine := flag.Bool("respect-max-line-length", false, "Respect max line length for soft wrapping.")
	abbrevs := flag.String("abbreviations", "", "Comma-separated list of custom abbreviations (e.g., \"No.,Fig.\").")

	flag.Parse()

	// Read all input from standard input.
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
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
	output, err := advent.Ventilate(string(input), cfg)
	if err != nil {
		log.Fatalf("Error ventilating text: %v", err)
	}

	// Write the result to standard output.
	fmt.Println(output)
}

