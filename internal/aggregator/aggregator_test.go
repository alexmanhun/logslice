package aggregator

import (
	"testing"
)

func makeTestEntries() []map[string]interface{} {
	return []map[string]interface{}{
		{"level": "error", "service": "auth"},
		{"level": "info", "service": "auth"},
		{"level": "error", "service": "api"},
		{"level": "warn", "service": "api"},
		{"level": "error", "service": "auth"},
		{"service": "db"},
	}
}

func TestCountBy_Level(t *testing.T) {
	entries := makeTestEntries()
	results := CountBy(entries, "level")

	if len(results) != 4 {
		t.Fatalf("expected 4 groups (error, info, warn, <missing>), got %d", len(results))
	}
	if results[0].Value != "error" || results[0].Count != 3 {
		t.Errorf("expected error=3, got %s=%d", results[0].Value, results[0].Count)
	}
}

func TestCountBy_MissingField(t *testing.T) {
	entries := makeTestEntries()
	results := CountBy(entries, "level")

	var missing *CountResult
	for i := range results {
		if results[i].Value == "<missing>" {
			missing = &results[i]
			break
		}
	}
	if missing == nil {
		t.Fatal("expected <missing> group")
	}
	if missing.Count != 1 {
		t.Errorf("expected <missing>=1, got %d", missing.Count)
	}
}

func TestUniqueValues_Service(t *testing.T) {
	entries := makeTestEntries()
	vals := UniqueValues(entries, "service")

	if len(vals) != 3 {
		t.Fatalf("expected 3 unique services, got %d", len(vals))
	}
	expected := []string{"api", "auth", "db"}
	for i, v := range vals {
		if v != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, v)
		}
	}
}

func TestUniqueValues_EmptyEntries(t *testing.T) {
	vals := UniqueValues([]map[string]interface{}{}, "level")
	if len(vals) != 0 {
		t.Errorf("expected empty result, got %v", vals)
	}
}
