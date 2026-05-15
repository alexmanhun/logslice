package sorter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/sorter"
)

func makeEntries() []map[string]any {
	return []map[string]any{
		{"ts": "2024-01-03", "level": "info", "msg": "third"},
		{"ts": "2024-01-01", "level": "error", "msg": "first"},
		{"ts": "2024-01-02", "level": "warn", "msg": "second"},
	}
}

func TestSortEntries_StringAscending(t *testing.T) {
	entries := makeEntries()
	if err := sorter.SortEntries(entries, "ts", sorter.Ascending); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"2024-01-01", "2024-01-02", "2024-01-03"}
	for i, e := range entries {
		if e["ts"] != expected[i] {
			t.Errorf("index %d: got %v, want %v", i, e["ts"], expected[i])
		}
	}
}

func TestSortEntries_StringDescending(t *testing.T) {
	entries := makeEntries()
	if err := sorter.SortEntries(entries, "ts", sorter.Descending); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"2024-01-03", "2024-01-02", "2024-01-01"}
	for i, e := range entries {
		if e["ts"] != expected[i] {
			t.Errorf("index %d: got %v, want %v", i, e["ts"], expected[i])
		}
	}
}

func TestSortEntries_NumericField(t *testing.T) {
	entries := []map[string]any{
		{"code": float64(300)},
		{"code": float64(100)},
		{"code": float64(200)},
	}
	if err := sorter.SortEntries(entries, "code", sorter.Ascending); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, want := range []float64{100, 200, 300} {
		if entries[i]["code"] != want {
			t.Errorf("index %d: got %v, want %v", i, entries[i]["code"], want)
		}
	}
}

func TestSortEntries_MissingFieldSinksToBottom(t *testing.T) {
	entries := []map[string]any{
		{"ts": "2024-01-02"},
		{"msg": "no ts here"},
		{"ts": "2024-01-01"},
	}
	if err := sorter.SortEntries(entries, "ts", sorter.Ascending); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := entries[2]["ts"]; ok {
		t.Errorf("expected entry without 'ts' to be last")
	}
}

func TestSortEntries_EmptyField(t *testing.T) {
	entries := makeEntries()
	err := sorter.SortEntries(entries, "", sorter.Ascending)
	if err == nil {
		t.Error("expected error for empty field, got nil")
	}
}

func TestSortEntries_TypeMismatch(t *testing.T) {
	entries := []map[string]any{
		{"val": "string"},
		{"val": float64(42)},
	}
	err := sorter.SortEntries(entries, "val", sorter.Ascending)
	if err == nil {
		t.Error("expected error for type mismatch, got nil")
	}
}
