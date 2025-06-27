# Advent Prose Ventilator

**Advent** is a Go library that performs structural ventilation of Markdown-formatted prose. It is designed to insert meaningful line breaks at natural sentence boundaries without altering the visual rendering of the document, making it ideal for improving the readability and diff-ability of prose in version control systems.

The library is markup-aware but markup-agnostic—it respects and preserves inline markup such as **EditML** and **CriticMarkup**, treating them as atomic spans during processing.

## Core Philosophy

* **Markdown-first**: Only processes prose blocks valid within Markdown.
* **Editor-driven**: Does not parse full documents with metadata or hybrid formats.
* **Markup-preserving**: Honors curly-braced inline markup but does not interpret it.
* **One-way formatting**: Advent does not provide "unventilation"; standard Markdown tools will reflow the output naturally.

---

## Installation

```sh
go get [github.com/verkaro/advent-prose](https://github.com/verkaro/advent-prose)
```

---

## Usage

The primary entrypoint to the library is the `advent.Ventilate` function.

### Basic Ventilation

To split a paragraph into one-sentence-per-line, use a simple `Config`.

```go
package main

import (
	"fmt"
	"log"

	"[github.com/verkaro/advent-prose](https://github.com/verkaro/advent-prose)"
)

func main() {
	input := "This is the first sentence. This is the second sentence. And this is the third."
	cfg := advent.Config{
		SentenceBreak: true,
	}

	output, err := advent.Ventilate(input, cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
```

**Output:**

```
This is the first sentence.
This is the second sentence.
And this is the third.
```

### Advanced Configuration

Advent provides several options to control its behavior.

#### Paragraph Spacing

You can control the spacing between blocks with `ParagraphSpacing`.

```go
input := "This is paragraph one.\n\nThis is paragraph two."
cfg := advent.Config{
    ParagraphSpacing: "blank-line",
}
// output will have a blank line between the paragraphs.
```

#### Custom Abbreviations

You can provide a custom list of abbreviations to prevent incorrect sentence breaks.

```go
input := "The item is No. 42. It is very important."
cfg := advent.Config{
    SentenceBreak: true,
    Abbreviations: map[string]bool{
        "No.": true,
    },
}

output, _ := advent.Ventilate(input, cfg)
// output:
// The item is No. 42.
// It is very important.
```

---

## Configuration

The `advent.Config` struct provides the following options:

| Field                  | Type              | Description                                                                                              |
| ---------------------- | ----------------- | -------------------------------------------------------------------------------------------------------- |
| `SentenceBreak`        | `bool`            | If `true`, inserts line breaks at the end of sentences.                                                  |
| `MaxLineLength`        | `int`             | A soft limit for line length when `RespectMaxLineLength` is `true`.                                      |
| `ParagraphSpacing`     | `string`          | Controls spacing between blocks. Can be `"single"` or `"blank-line"`.                                      |
| `RespectMaxLineLength` | `bool`            | If `true` and `SentenceBreak` is `false`, wraps lines at word boundaries.                                |
| `Abbreviations`        | `map[string]bool` | A custom map of abbreviations to ignore. If `nil`, a default English list is used.                       |


