# Advent Spec v0.3 — Markdown Prose Ventilator Library

## Overview

Advent is a Go library for structurally ventilating Markdown-formatted prose. It inserts line breaks at natural sentence boundaries for improved readability, diff tracking, and editing. It is compatible with inline markup like **EditML** and **CriticMarkup**, and preserves the integrity of Markdown documents.

While Advent operates on raw Markdown strings, it is designed to align with the [Goldmark](https://github.com/yuin/goldmark) parser model. This ensures that the output remains semantically valid and visually unchanged when rendered through Goldmark or any compliant Markdown renderer.

> **Note:** Advent only processes prose blocks. Structural blocks (headings, lists, code, etc.) are preserved without modification.

---

## Highlights

* ✨ **Markdown-aware:** Understands and preserves Markdown structure
* 🧠 **Markup-smart:** Protects `{...}` inline markup (EditML, CriticMarkup)
* 📐 **Configurable:** Supports sentence-aware or line-length-based reflow
* 🔄 **Idempotent:** Repeated runs will not mangle spacing or breaks
* 🧪 **Testable:** Ships with a full test suite and golden fixtures

---

## Public API

### `type Config`

Defines global configuration for sentence detection and formatting.

```go
type Config struct {
    SentenceBreak           bool   // Break after sentence-ending punctuation
    MaxLineLength           int    // Soft max characters per line
    ParagraphSpacing        string // "none", "single", or "blank-line"
    RespectMaxLineLength    bool   // Wrap long lines if no sentence break
}
```

### `func Ventilate(input string, cfg Config) (string, error)`

Performs sentence-aware reflow on Markdown-formatted prose paragraphs.

### `func IsVentilated(input string) bool`

(Optional) Detects whether input already follows Advent-style ventilation.

---

## Input Requirements

* Valid Markdown prose (paragraphs, emphasis, links, etc.)
* May include inline `{...}` markup, such as:

  * `{+add+}`, `{-remove-}`, `{=highlight=}`, `{>>comment<<}`
  * Nested punctuation, movement tags, or combinations
* Should not contain structural or config-driven syntax (e.g. `.biff`)

---

## Sentence Detection Behavior

* Breaks after `.`, `!`, or `?` followed by space or newline
* Uses a masked version of the string to ignore punctuation in markup
* No breaks inserted inside `{...}` spans
* Sentence-ending punctuation inside quotes or emphasis is respected
* Allowlist of abbreviations prevents false sentence breaks

---

## Markdown-Aware Formatting

Advent respects Markdown syntax rules:

| Structure   | Behavior           |
| ----------- | ------------------ |
| Paragraphs  | Ventilated         |
| Headings    | Passed unchanged   |
| Lists       | Passed unchanged   |
| Blockquotes | Passed unchanged   |
| Code blocks | Preserved verbatim |

---

## Configuration Behavior

* If `SentenceBreak: true`, never break mid-sentence—even if too long
* If `RespectMaxLineLength: true`, wrap at word boundaries (unless SentenceBreak is on)
* Paragraph spacing can be:

  * `"none"` — flush paragraph joins
  * `"single"` — single newline
  * `"blank-line"` — extra blank line

---

## Error Conditions

| Situation                        | Error? |
| -------------------------------- | ------ |
| Malformed `{+add`                | ✅      |
| Empty string input               | ❌      |
| Invalid nested braces `{=x{y}=}` | ✅      |

---

## Newline Normalization

* Trims leading/trailing whitespace
* Collapses multiple blank lines to one
* Ensures output is consistent with configuration

---

## Testing & Verification

* ✅ Golden tests for all edge cases
* ✅ Table-driven tests for configuration variations
* ✅ Idempotency checks
* ✅ Malformed markup triggers error conditions
* ✅ Markdown output is human-readable and stable for version control

> See `test_advent_expected.md` for examples.

---

## Origin

Advent was developed collaboratively:

* Specification authored and iterated with **ChatGPT**
* Initial implementation by **Gemini**

Together, they rapidly refined the tool into an extensible, testable, and Markdown-compliant prose processor.

---

## License

MIT. See `LICENSE` file.

---



