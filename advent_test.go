package advent

import (
	"errors"
	"testing"
)

// TestVentilate runs a table-driven test to cover all the edge cases
// and behaviors defined in the Advent specification.
func TestVentilate(t *testing.T) {
	testCases := []struct {
		name   string
		cfg    Config
		input  string
		want   string
		err    error
	}{
		// --- Existing test cases ---

		{
			name:  "Simple paragraph",
			cfg:   Config{SentenceBreak: true},
			input: "This is a sentence. Here is another.",
			want:  "This is a sentence.\nHere is another.",
			err:   nil,
		},
		{
			name:  "EditML with punctuation",
			cfg:   Config{SentenceBreak: true},
			input: "This is the beginning.{+ And it continues.+} But is it?",
			want:  "This is the beginning.{+ And it continues.+}\nBut is it?",
			err:   nil,
		},
		{
			name:  "CriticMarkup highlight and comment",
			cfg:   Config{SentenceBreak: true},
			input: "{=There was a dog=}{>this is passive voice!<} in the corner of the room. It might have been rabid{-?-}{+.+}",
			want:  "{=There was a dog=}{>this is passive voice!<} in the corner of the room.\nIt might have been rabid{-?-}{+.+}",
			err:   nil,
		},
		{
			name:  "Emphasis includes terminal punctuation",
			cfg:   Config{SentenceBreak: true},
			input: "I like *strong tea.* It helps me write.",
			want:  "I like *strong tea.*\nIt helps me write.",
			err:   nil,
		},
		{
			name:  "Sentence-ending punctuation within quotes",
			cfg:   Config{SentenceBreak: true},
			input: `"I lived without much I missed!" she said. It was true.`,
			want:  `"I lived without much I missed!" she said.` + "\n" + `It was true.`,
			err:   nil,
		},
		{
			name:  "Markdown Blockquote With Sentence Punctuation (No-op)",
			cfg:   Config{SentenceBreak: true},
			input: "> This is a line. It should not break.",
			want:  "> This is a line. It should not break.",
			err:   nil,
		},
		{
			name:  "Fenced Code Block With Sentences (No-op)",
			cfg:   Config{SentenceBreak: true},
			input: "```\nHere is some code. It does not need ventilation.\nprint(\"Hello world.\")\n```",
			want:  "```\nHere is some code. It does not need ventilation.\nprint(\"Hello world.\")\n```",
			err:   nil,
		},
		{
			name:  "Markdown Heading With Sentence-Like Content (No-op)",
			cfg:   Config{SentenceBreak: true},
			input: "# My Story. A Lasting Relationship.",
			want:  "# My Story. A Lasting Relationship.",
			err:   nil,
		},
		{
			name:  "Markdown List With Punctuation (No-op)",
			cfg:   Config{SentenceBreak: true},
			input: "- Found. A dog.",
			want:  "- Found. A dog.",
			err:   nil,
		},
		{
			name:  "Abbreviation with default list",
			cfg:   Config{SentenceBreak: true},
			input: "I went to see Mr. Smith. He was home.",
			want:  "I went to see Mr. Smith.\nHe was home.",
			err:   nil,
		},
		{
			name: "Config: MaxLineLength without SentenceBreak",
			cfg: Config{
				SentenceBreak:        false,
				MaxLineLength:        20,
				RespectMaxLineLength: true,
			},
			input: "This is a long sentence that should wrap gently.",
			want:  "This is a long\nsentence that should\nwrap gently.",
			err:   nil,
		},
		{
			name:  "Error: Unterminated Markup",
			cfg:   Config{SentenceBreak: true},
			input: "Something went wrong {+unfinished",
			want:  "", // Expect no output on error
			err:   errors.New("unterminated inline markup span"),
		},
		{
			name:  "Markdown Link Wrapping",
			cfg:   Config{SentenceBreak: true},
			input: "You can [learn more about our project](https://example.com). It is worth a visit.",
			want:  "You can [learn more about our project](https://example.com).\nIt is worth a visit.",
			err:   nil,
		},
		{
			name:  "Strong Emphasis Ending Sentence",
			cfg:   Config{SentenceBreak: true},
			input: "This is **very important.** Please take note.",
			want:  "This is **very important.**\nPlease take note.",
			err:   nil,
		},
		{
			name:  "CriticMarkup with nested punctuation",
			cfg:   Config{SentenceBreak: true},
			input: "The trial was unfair{>>needs clarification!<<}. Justice must be served.",
			want:  "The trial was unfair{>>needs clarification!<<}.\nJustice must be served.",
			err:   nil,
		},
		{
			name: "ParagraphSpacing: single",
			cfg: Config{
				SentenceBreak:    true,
				ParagraphSpacing: "single",
			},
			input: "Paragraph one. It has two sentences.\n\nParagraph two. Also two sentences.",
			want:  "Paragraph one.\nIt has two sentences.\nParagraph two.\nAlso two sentences.",
			err:   nil,
		},
		{
			name: "ParagraphSpacing: blank-line",
			cfg: Config{
				SentenceBreak:    true,
				ParagraphSpacing: "blank-line",
			},
			input: "Paragraph one.\n\nParagraph two.",
			want:  "Paragraph one.\n\nParagraph two.",
			err:   nil,
		},
		{
			name: "ParagraphSpacing: blank-line between prose and heading",
			cfg: Config{
				ParagraphSpacing: "blank-line",
			},
			input: "This is prose.\n# This is a heading",
			want:  "This is prose.\n\n# This is a heading",
			err:   nil,
		},

		// --- New test for configurable abbreviations ---
		{
			name: "Abbreviation with custom list",
			cfg: Config{
				SentenceBreak: true,
				Abbreviations: map[string]bool{
					"No.": true, // Custom abbreviation
				},
			},
			input: "Item No. 42 is important. The default Mr. Smith is not an abbreviation here.",
			// Expect "No. 42" to be preserved, but "Mr." to be split.
			want:  "Item No. 42 is important.\nThe default Mr.\nSmith is not an abbreviation here.",
			err:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Ventilate(tc.input, tc.cfg)

			if (err != nil && tc.err == nil) || (err == nil && tc.err != nil) || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
				t.Errorf("Ventilate() error = %v, wantErr %v", err, tc.err)
				return
			}

			if got != tc.want {
				t.Errorf("Ventilate() got = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestIsVentilated tests the helper function for detecting ventilation.
func TestIsVentilated(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Already ventilated text",
			input: "This is a sentence.\nHere is another one.",
			want:  true,
		},
		{
			name:  "Unventilated text",
			input: "This is a sentence. Here is another one.",
			want:  false,
		},
		{
			name:  "Empty string",
			input: "",
			want:  true, // An empty string is trivially ventilated.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsVentilated(tc.input); got != tc.want {
				t.Errorf("IsVentilated() = %v, want %v", got, tc.want)
			}
		})
	}
}

