package pipeline_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/yourorg/logslice/internal/pipeline"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return f.Name()
}

func TestRun_BasicFilter(t *testing.T) {
	file := writeTempLog(t, []string{
		`{"level":"info","msg":"started"}`,
		`{"level":"error","msg":"boom"}`,
		`{"level":"info","msg":"done"}`,
	})

	var buf bytes.Buffer
	err := pipeline.Run(pipeline.Config{
		Files:  []string{file},
		Filter: "level=error",
		Format: "json",
	}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if entry["level"] != "error" {
		t.Errorf("expected level=error, got %v", entry["level"])
	}
}

func TestRun_Limit(t *testing.T) {
	file := writeTempLog(t, []string{
		`{"level":"info","n":1}`,
		`{"level":"info","n":2}`,
		`{"level":"info","n":3}`,
	})

	var buf bytes.Buffer
	_ = pipeline.Run(pipeline.Config{
		Files:  []string{file},
		Limit:  2,
		Format: "json",
	}, &buf)

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestRun_CountBy(t *testing.T) {
	file := writeTempLog(t, []string{
		`{"level":"info"}`,
		`{"level":"info"}`,
		`{"level":"error"}`,
	})

	var buf bytes.Buffer
	_ = pipeline.Run(pipeline.Config{
		Files:   []string{file},
		CountBy: "level",
	}, &buf)

	if buf.Len() == 0 {
		t.Error("expected non-empty aggregation output")
	}
}

func TestRun_InvalidFilter(t *testing.T) {
	var buf bytes.Buffer
	err := pipeline.Run(pipeline.Config{
		Files:  []string{},
		Filter: "!!!",
	}, &buf)
	if err == nil {
		t.Error("expected error for invalid filter")
	}
}
