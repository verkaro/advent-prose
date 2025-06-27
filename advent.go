// Package advent provides tools for structural ventilation of Markdown-formatted prose.
// It inserts line breaks at natural sentence boundaries to improve readability and
// version control diffs without altering the final rendered output.
package advent

import (
	"errors"
	"strings"
	"unicode"
)

// Config holds the configuration options for the Ventilate function.
// It allows callers to control how Markdown prose is ventilated.
type Config struct {
	// SentenceBreak, when true, instructs the ventilator to insert line breaks
	// at the end of sentences.
	SentenceBreak bool

	// MaxLineLength is an optional soft maximum number of characters per line.
	// If RespectMaxLineLength is false, this is ignored.
	MaxLineLength int

	// ParagraphSpacing defines the style of spacing between paragraphs.
	// Valid values: "none", "single", "blank-line".
	ParagraphSpacing string

	// RespectMaxLineLength, if true, soft-wraps long lines at word boundaries,
	// as long as SentenceBreak is false.
	RespectMaxLineLength bool

	// Abbreviations is a map of custom abbreviations to prevent sentence splitting.
	// If nil, a default list of English abbreviations is used.
	Abbreviations map[string]bool
}

var (
	errUnterminatedMarkup = errors.New("unterminated inline markup span")

	// defaultAbbreviations is a set of common English abbreviations that
	// should not be treated as the end of a sentence.
	defaultAbbreviations = map[string]bool{
		"Mr.": true, "Mrs.": true, "Ms.": true, "Dr.": true, "Prof.": true,
		"e.g.": true, "i.e.": true, "etc.": true, "St.": true,
	}
)

// Ventilate performs sentence-aware reflow on Markdown-formatted prose paragraphs.
func Ventilate(input string, cfg Config) (string, error) {
	if input == "" {
		return "", nil
	}

	if err := checkUnterminatedMarkup(input); err != nil {
		return "", err
	}

	// Preserve trailing newline information.
	hasTrailingNewline := strings.HasSuffix(input, "\n") || strings.HasSuffix(input, "\r\n")

	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	var processedBlocks []string
	var blockBuffer []string

	// processBlock is a helper to process the buffered lines as a single block.
	processBlock := func() error {
		if len(blockBuffer) > 0 {
			processed, err := ventilateBlock(blockBuffer, cfg)
			if err != nil {
				return err
			}
			processedBlocks = append(processedBlocks, processed)
			blockBuffer = nil // Reset buffer
		}
		return nil
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			// A blank line acts as a block separator.
			if err := processBlock(); err != nil {
				return "", err
			}
		} else {
			blockBuffer = append(blockBuffer, line)
		}
	}

	// Process the last block if the file doesn't end with a blank line.
	if err := processBlock(); err != nil {
		return "", err
	}

	// Join all processed blocks with a standard double newline.
	output := strings.Join(processedBlocks, "\n\n")

	// Restore trailing newline if it existed in the original input.
	if hasTrailingNewline && !strings.HasSuffix(output, "\n") {
		output += "\n"
	}


	return output, nil
}

// ventilateBlock determines the block type and processes it accordingly.
func ventilateBlock(blockLines []string, cfg Config) (string, error) {
	// A block is non-prose if its first line indicates a non-prose type.
	firstLineTrimmed := strings.TrimSpace(blockLines[0])

	if isNonProseBlock(firstLineTrimmed) || strings.HasPrefix(firstLineTrimmed, "```") {
		// Non-prose blocks are passed through verbatim, preserving internal newlines.
		return strings.Join(blockLines, "\n"), nil
	}

	// Join lines with care, preserving meaningful breaks like those after a colon.
	var paraBuilder strings.Builder
	for i, line := range blockLines {
		paraBuilder.WriteString(line)
		if i < len(blockLines)-1 {
			// If a line ends with a colon, preserve the newline.
			// Otherwise, join with a space to merge wrapped lines.
			if strings.HasSuffix(strings.TrimSpace(line), ":") {
				paraBuilder.WriteString("\n")
			} else {
				paraBuilder.WriteString(" ")
			}
		}
	}

	return ventilateParagraph(paraBuilder.String(), cfg)
}

// ventilateParagraph handles the core logic for a single prose paragraph.
func ventilateParagraph(p string, cfg Config) (string, error) {
	if cfg.SentenceBreak {
		return ventilateBySentence(p, cfg)
	}
	if cfg.RespectMaxLineLength && cfg.MaxLineLength > 0 {
		return ventilateByLineLength(p, cfg.MaxLineLength)
	}
	return p, nil
}

// ventilateBySentence processes a paragraph, inserting newlines after sentences.
func ventilateBySentence(p string, cfg Config) (string, error) {
	var result strings.Builder
	lastBreak := 0
	i := 0

	abbreviations := cfg.Abbreviations
	if abbreviations == nil {
		abbreviations = defaultAbbreviations
	}

	for i < len(p) {
		if strings.HasPrefix(p[i:], "...") {
			i += 3
			continue
		}

		if p[i] == '{' {
			end, ok := findMarkupEnd(p, i)
			if !ok {
				return "", errUnterminatedMarkup
			}
			i = end
			j := i + 1
			for j < len(p) && unicode.IsSpace(rune(p[j])) {
				j++
			}
			if j >= len(p) || unicode.IsUpper(rune(p[j])) {
				result.WriteString(p[lastBreak : i+1])
				result.WriteRune('\n')
				lastBreak = j
				i = j
				continue
			}
			i++
			continue
		}

		char := p[i]
		if char == '.' || char == '!' || char == '?' {
			wordStart := strings.LastIndexAny(p[:i], " \n")
			if wordStart == -1 {
				wordStart = 0
			} else {
				wordStart++
			}
			word := p[wordStart : i+1]
			if abbreviations[word] {
				i++
				continue
			}

			if (char == '!' || char == '?') && i+1 < len(p) && p[i+1] == '"' {
				j := i + 2
				for j < len(p) && unicode.IsSpace(rune(p[j])) {
					j++
				}
				if j < len(p) && unicode.IsLower(rune(p[j])) {
					i++
					continue
				}
			}

			j := i + 1
			for j < len(p) {
				if strings.ContainsRune("*}_)]}\"'", rune(p[j])) {
					j++
				} else {
					break
				}
			}

			if j >= len(p) || unicode.IsSpace(rune(p[j])) {
				result.WriteString(p[lastBreak:j])
				result.WriteRune('\n')
				lastBreak = j
				for lastBreak < len(p) && unicode.IsSpace(rune(p[lastBreak])) {
					lastBreak++
				}
				i = lastBreak
				continue
			}
		}
		i++
	}

	if lastBreak < len(p) {
		result.WriteString(p[lastBreak:])
	}
	return strings.TrimSuffix(result.String(), "\n"), nil
}

func findMarkupEnd(s string, start int) (int, bool) {
	level := 1
	for i := start + 1; i < len(s); i++ {
		if s[i] == '{' {
			level++
		} else if s[i] == '}' {
			level--
			if level == 0 {
				return i, true
			}
		}
	}
	return -1, false
}

func ventilateByLineLength(p string, maxLen int) (string, error) {
	var result strings.Builder
	words := strings.Fields(p)
	if len(words) == 0 {
		return "", nil
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) > maxLen {
			result.WriteString(currentLine)
			result.WriteRune('\n')
			currentLine = word
		} else {
			currentLine += " " + word
		}
	}
	result.WriteString(currentLine)

	return result.String(), nil
}

func checkUnterminatedMarkup(s string) error {
	level := 0
	for _, r := range s {
		if r == '{' {
			level++
		} else if r == '}' {
			if level > 0 {
				level--
			} else {
				return errUnterminatedMarkup
			}
		}
	}
	if level != 0 {
		return errUnterminatedMarkup
	}
	return nil
}

// isNonProseBlock determines if a line marks the beginning of a block
// that should be passed through without ventilation.
func isNonProseBlock(line string) bool {
	if line == "" {
		return false
	}
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ">") || strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "---") {
		return true
	}
	if i := strings.Index(line, ". "); i > 0 {
		numPart := line[:i]
		if len(numPart) == 0 {
			return false
		}
		isAllDigits := true
		for _, r := range numPart {
			if !unicode.IsDigit(r) {
				isAllDigits = false
				break
			}
		}
		if isAllDigits {
			return true
		}
	}
	return false
}

// IsVentilated detects whether the input string already follows the Advent reflow style.
func IsVentilated(input string) bool {
	if input == "" {
		return true
	}
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		for _, punc := range []string{". ", "! ", "? "} {
			if idx := strings.Index(line, punc); idx != -1 {
				if len(strings.TrimSpace(line[idx+2:])) > 0 {
					return false
				}
			}
		}
	}
	return true
}

