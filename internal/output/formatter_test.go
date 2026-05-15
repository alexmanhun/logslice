package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func sampleEntry() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello world",
		"ts":      1700000000,
	}
}

func TestFormatter_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["level"] != "info" {
		t.Errorf("expected level=info, got %v", got["level"])
	}
}

func TestFormatter_WritePretty(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatPretty)

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\n") {
		t.Error("expected pretty output to contain newlines")
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &got); err != nil {
		t.Fatalf("pretty output is not valid JSON: %v", err)
	}
}

func TestFormatter_WriteText(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected text output to contain 'level=info', got: %s", out)
	}
	if !strings.Contains(out, "message=hello world") {
		t.Errorf("expected text output to contain 'message=hello world', got: %s", out)
	}
}

func TestFormatter_DefaultIsJSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, "")

	if err := f.Write(sampleEntry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("default output is not valid JSON: %v", err)
	}
}
