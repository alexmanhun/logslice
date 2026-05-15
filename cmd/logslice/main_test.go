package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, lines []map[string]interface{}) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, l := range lines {
		if err := enc.Encode(l); err != nil {
			t.Fatalf("encode line: %v", err)
		}
	}
	return f.Name()
}

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "logslice")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("could not build binary: %v\n%s", err, out)
	}
	return bin
}

func TestMain_NoArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "-help")
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "Usage:") {
		t.Errorf("expected usage output, got: %s", out)
	}
}

func TestMain_FilterByLevel(t *testing.T) {
	bin := buildBinary(t)
	file := writeTempLog(t, []map[string]interface{}{
		{"level": "error", "msg": "boom"},
		{"level": "info", "msg": "ok"},
	})

	cmd := exec.Command(bin, "-query", "level=error", "-format", "json", file)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %s", len(lines), out)
	}
	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if entry["level"] != "error" {
		t.Errorf("expected level=error, got %v", entry["level"])
	}
}

func TestMain_InvalidQuery(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "-query", "!!!", "somefile.log")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit for invalid query")
	}
	if !strings.Contains(string(out), "invalid query") {
		t.Errorf("expected 'invalid query' in stderr, got: %s", out)
	}
}
