# Advent Spec v0.3 — Markdown Prose Ventilator Library

## Overview

Advent is designed to be used alongside the [Goldmark](https://github.com/yuin/goldmark) Markdown parser in Go. While Advent operates on raw Markdown strings, its architecture and behavior are guided by Goldmark's block and inline parsing model. This ensures compatibility and future extensibility for projects already using Goldmark for Markdown rendering or analysis.
**Advent** is a Go library that performs structural ventilation of Markdown-formatted prose for editing, reviewing, and version control purposes. It is designed to insert meaningful line breaks at natural sentence boundaries without altering the visual rendering of the document.

Advent is **markup-aware but markup-agnostic** — it respects and preserves inline markup such as **EditML** and **CriticMarkup**, treating them as atomic spans during processing.

---

## Core Philosophy

* **Markdown-first**: Only processes prose blocks valid within Markdown.
* **Editor-driven**: Does not parse full documents with metadata, game config syntax, or hybrid formats like `.biff`.
* **Markup-preserving**: Honors curly-braced inline markup but does not interpret it.
* **One-way formatting**: Advent does not provide "unventilation"; tools like Goldmark will reflow output naturally.

---

## Public API

### `type Config`

Defines global configuration for sentence detection and formatting.

```go
type Config struct {
    // Whether to insert line breaks after sentence-ending punctuation.
    SentenceBreak bool

    // Optional soft maximum number of characters per line.
    // If RespectMaxLineLength is false, this is ignored.
    MaxLineLength int

    // Controls inter-paragraph spacing.
    // Valid values: "none", "single", "blank-line"
    ParagraphSpacing string

    // If true, soft-wrap long lines at word boundaries,
    // as long as SentenceBreak is false.
    RespectMaxLineLength bool
}
```

### `func Ventilate(input string, cfg Config) (string, error)`

Performs sentence-aware reflow on Markdown-formatted prose paragraphs.

```go
func Ventilate(input string, cfg Config) (string, error)
```

* Ignores headings, code blocks, lists, etc.
* Preserves inline markup spans (e.g. `{+add+}`, `{=highlight=}`)
* Returns the transformed Markdown with appropriate line breaks

### `func IsVentilated(input string) bool`

(Optional) Detects whether input already follows the Advent reflow style.

```go
func IsVentilated(input string) bool
```

* Could be used in editor plugins or tests to avoid double-formatting

---

## Input Requirements

* Valid Markdown prose only (paragraphs, headers, emphasis, etc.)
* May include inline markup using `{...}` braces, such as:

  * EditML: `{+add+}`, `{-remove-}`, `{=highlight=}`, `{move~text~id}`
  * CriticMarkup: `{++add++}`, `{--delete--}`, `{==highlight==}`, `{>>comment<<}`
* Should not include structural logic or non-Markdown syntaxes (e.g., `.biff` scenes)

---

## Behavior

Advent is informed by the structural expectations of Goldmark. If integrated into a Goldmark-based toolchain, it can operate on pre-parsed Markdown blocks (e.g. paragraphs or inline text spans) and return ventilated Markdown content that respects formatting boundaries. While Advent itself does not require Goldmark, its output is intentionally structured to preserve Markdown integrity in Goldmark’s rendering pipeline.

### Structural Block Awareness

* Advent does **not** process or ventilate non-prose Markdown blocks, including:

  * Headings (e.g., lines beginning with `#`)
  * List items (lines beginning with `-`, `*`, or numbered lists like `1.`)
  * Code blocks, blockquotes, tables, or other container blocks
* Sentence punctuation found within these blocks is preserved verbatim, and no line breaks are inserted
* Only **prose paragraphs** are processed for sentence ventilation

### ✅ Sentence Break Handling

* Breaks inserted **after sentence-ending punctuation** (`.`, `!`, `?`) followed by whitespace or newline
* Sentence detection uses a **masked representation** of the input:

  * All `{...}` spans are replaced with equal-length placeholders during sentence detection
  * This ensures punctuation within markup is ignored
  * The original text is then broken at safe boundaries
* Ellipses (`...`) are not treated as sentence breaks unless context clearly indicates end of sentence
* Quoted punctuation (e.g. `"...!" she said.`) is treated as sentence-ending
* Abbreviations such as `e.g.`, `Mr.`, `St.` are preserved using a configurable allowlist to prevent premature breaks

### ✅ Markup Protection

* No line breaks are ever inserted **inside** a `{...}` span
* Inline markup is preserved exactly as input

### ✅ Paragraph Formatting

* Newlines are inserted between sentences
* Existing newlines are normalized based on `cfg`
* Paragraph-level Markdown structure (e.g. emphasis, strong, links) is preserved
* `ParagraphSpacing` may be:

  * "none" (default): no extra spacing
  * "single": single newline between paragraphs
  * "blank-line": adds a full blank line between paragraphs

---

## Configuration Interactions

* If `SentenceBreak` is `true`, Advent only breaks between **complete sentences** — never mid-sentence — even if `MaxLineLength` is exceeded
* If `RespectMaxLineLength` is `true`, and `SentenceBreak` is `false`, Advent may soft-wrap long lines at the nearest safe word boundary, excluding `{...}` spans
* If both are `false`, Advent does not perform line breaking

---

## Error Conditions

* The `Ventilate` function will return an error for:

  * Unterminated inline markup spans (e.g. `{+unclosed`)
  * Unrecognized or malformed `{...}` syntax if `cfg.SentenceBreak` is `true` and masking fails
  * Nil or empty input will return an empty result and no error

---

## Normalization of Newlines

* Leading and trailing whitespace is trimmed from each line
* Multiple blank lines between paragraphs are collapsed to the spacing defined in `cfg.ParagraphSpacing`
* If input is already ventilated, Advent will normalize to its configured format (idempotent behavior encouraged)

---

## Edge Case Examples

### 1. Simple paragraph

**Input:**

```
This is a sentence. Here is another.
```

**Output:**

```
This is a sentence.
Here is another.
```

### 2. EditML with punctuation

**Input:**

```
This is the beginning.{+ And it continues.+} But is it?
```

**Output:**

```
This is the beginning.{+ And it continues.+}
But is it?
```

### 3. CriticMarkup highlight and comment

**Input:**

```
{=There was a dog=}{>this is passive voice!<} in the corner of the room. It might have been rabid{-?-}{+.+}
```

**Output:**

```
{=There was a dog=}{>this is passive voice!<} in the corner of the room.
It might have been rabid{-?-}{+.+}
```

### 4. Sentences near link or emphasis

**Input:**

```
I like *strong tea*. It helps me write.
```

**Output:**

```
I like *strong tea*.
It helps me write.
```

### 5. Emphasis includes terminal punctuation

**Input:**

```
I like *strong tea.* It helps me write.
```

**Output:**

```
I like *strong tea.*
It helps me write.
```

### 6. Ellipsis mid-sentence

**Input:**

```
He looked up... the stars were already fading. She whispered.
```

**Output:**

```
He looked up... the stars were already fading.
She whispered.
```

### 7. Sentence-ending punctuation within quotes

**Input:**

```
"I lived without much I missed!" she said. It was true.
```

**Output:**

```
"I lived without much I missed!" she said.
It was true.
```

---

## Testing & Verification

* Advent should be accompanied by a set of automated unit tests that:

  * Verify sentence detection and line break insertion for all documented edge cases
  * Assert preservation of inline markup spans
  * Confirm normalization behavior according to config settings
* Expected outputs for edge cases should be encoded as fixtures or table-driven test cases
* Behavior should be deterministic and idempotent
* Markdown inputs with malformed markup should trigger defined error conditions

---

## Test Case Appendix

### ✅ Markdown Blockquote With Sentence Punctuation

**Input:**

```
> This is a line. It should not break.
```

**Output:**

```
> This is a line. It should not break.
```

✅ Blockquotes are passed through unchanged regardless of sentence punctuation.

### ✅ Fenced Code Block With Sentences

**Input:**

````
```
Here is some code. It does not need ventilation.
print("Hello world.")
```
````

**Output:**

````
```
Here is some code. It does not need ventilation.
print("Hello world.")
```
````

✅ Code blocks are preserved entirely without modification.

### ✅ Markdown Heading With Sentence-Like Content

**Input:**

```
# My Story. A Lasting Relationship.
```

**Output:**

```
# My Story. A Lasting Relationship.
```

✅ Headings are not ventilated regardless of sentence punctuation.

### ✅ Markdown List With Punctuation

**Input:**

```
- Found. A dog.
```

**Output:**

```
- Found. A dog.
```

✅ List items retain inline punctuation without breaking across lines.

### ✅ Abbreviation

**Input:**

```
I went to see Mr. Smith. He was home.
```

**Output:**

```
I went to see Mr. Smith.
He was home.
```

### ✅ Config: MaxLineLength without SentenceBreak

**Config:**

```
SentenceBreak: false
MaxLineLength: 20
RespectMaxLineLength: true
```

**Input:**

```
This is a long sentence that should wrap gently.
```

**Output:**

```
This is a long
sentence that should
wrap gently.
```

### ✅ Error: Unterminated Markup

**Input:**

```
Something went wrong {+unfinished
```

**Expect Error:** true

### ✅ Markdown Link Wrapping

**Input:**

```
You can [learn more about our project](https://example.com). It is worth a visit.
```

**Output:**

```
You can [learn more about our project](https://example.com).
It is worth a visit.
```

### ✅ Emphasis Mid-Sentence

**Input:**

```
The *quick brown fox* jumps over the lazy dog. Again.
```

**Output:**

```
The *quick brown fox* jumps over the lazy dog.
Again.
```

### ✅ Strong Emphasis Ending Sentence

**Input:**

```
This is **very important.** Please take note.
```

**Output:**

```
This is **very important.**
Please take note.
```

### ✅ Link Inside Sentence Without Line Break

**Input:**

```
Refer to the [manual](doc.md) for more. Additional info follows.
```

**Output:**

```
Refer to the [manual](doc.md) for more.
Additional info follows.
```

### ✅ CriticMarkup with nested punctuation

**Input:**

```
The trial was unfair{>>needs clarification!<<}. Justice must be served.
```

**Output:**

```
The trial was unfair{>>needs clarification!<<}.
Justice must be served.
```

---

## Not in Scope

* `.biff` or any other hybrid config/markup formats
* Rendering or diffing tools
* Semantic analysis of prose
* Markdown-to-HTML conversion

---

## Future Considerations

* Optional Markdown AST traversal for smarter spacing
* Configurable breakpoints (e.g. clause-level)
* Plugin support for editor integrations (Argos, Soma, etc.)

