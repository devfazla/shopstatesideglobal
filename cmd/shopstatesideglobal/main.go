// Command shopstatesideglobal is the CLI entry point for the Stateside Global
// tooling. Today it exposes the HTML-to-Markdown cleaner; more subcommands and
// flags will be added as the codebase grows.
//
// Usage:
//
//	shopstatesideglobal --clean-html "<p>Hello <b>world</b></p>"
//	shopstatesideglobal --clean-html-file page.html
//	cat page.html | shopstatesideglobal --clean-html -
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/devfazla/shopstatesideglobal/internal/markdown"
)

func main() {
	var (
		cleanHTML     = flag.String("clean-html", "", "convert the given HTML string to clean Markdown and print it")
		cleanHTMLFile = flag.String("clean-html-file", "", "read HTML from a file and print clean Markdown")
	)
	flag.Parse()

	out, err := run(*cleanHTML, *cleanHTMLFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if out != "" {
		fmt.Println(out)
	}
}

// run resolves the HTML input from the flags / stdin and returns the cleaned
// Markdown. Keeping this separate from main makes the behaviour testable.
func run(cleanHTML, cleanHTMLFile string) (string, error) {
	switch {
	case strings.TrimSpace(cleanHTMLFile) != "":
		data, err := os.ReadFile(cleanHTMLFile)
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", cleanHTMLFile, err)
		}
		return markdown.CleanHTML(string(data)), nil

	case cleanHTML == "-":
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return markdown.CleanHTML(string(data)), nil

	case strings.TrimSpace(cleanHTML) != "":
		return markdown.CleanHTML(cleanHTML), nil

	default:
		// No explicit input: fall back to piped stdin if present.
		if hasStdin() {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", fmt.Errorf("reading stdin: %w", err)
			}
			return markdown.CleanHTML(string(data)), nil
		}
	}

	return "", nil
}

// hasStdin reports whether stdin is a pipe/file with data (not an interactive
// terminal). On Windows this treats stdin as available when it is not a char
// device.
func hasStdin() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}
