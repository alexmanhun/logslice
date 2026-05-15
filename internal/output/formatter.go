// Package output provides formatting utilities for structured log output.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format controls how log entries are rendered.
type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
	FormatText   Format = "text"
)

// Formatter writes log entries to an output stream.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a new Formatter writing to w with the given format.
func New(w io.Writer, format Format) *Formatter {
	return &Formatter{w: w, format: format}
}

// Write renders a single log entry to the underlying writer.
func (f *Formatter) Write(entry map[string]any) error {
	switch f.format {
	case FormatPretty:
		return f.writePretty(entry)
	case FormatText:
		return f.writeText(entry)
	default:
		return f.writeJSON(entry)
	}
}

func (f *Formatter) writeJSON(entry map[string]any) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(f.w, string(b))
	return err
}

func (f *Formatter) writePretty(entry map[string]any) error {
	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(f.w, string(b))
	return err
}

func (f *Formatter) writeText(entry map[string]any) error {
	keys := make([]string, 0, len(entry))
	for k := range entry {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(entry))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, entry[k]))
	}
	_, err := fmt.Fprintln(f.w, strings.Join(parts, " "))
	return err
}
