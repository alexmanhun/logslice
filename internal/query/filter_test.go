package query

import (
	"testing"
)

func TestParse_ValidQueries(t *testing.T) {
	tests := []struct {
		input    string
		wantLen  int
		wantFirst Filter
	}{
		{"level=error", 1, Filter{"level", OpEq, "error"}},
		{"level!=debug", 1, Filter{"level", OpNeq, "debug"}},
		{"msg~timeout", 1, Filter{"msg", OpContains, "timeout"}},
		{"status>200", 1, Filter{"status", OpGt, "200"}},
		{"status<500", 1, Filter{"status", OpLt, "500"}},
		{"level=error msg~fail", 2, Filter{"level", OpEq, "error"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d filters, got %d", tt.wantLen, len(got))
			}
			if got[0] != tt.wantFirst {
				t.Errorf("first filter: got %+v, want %+v", got[0], tt.wantFirst)
			}
		})
	}
}

func TestParse_EmptyQuery(t *testing.T) {
	filters, err := Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Errorf("expected no filters, got %d", len(filters))
	}
}

func TestParse_InvalidToken(t *testing.T) {
	_, err := Parse("invalidtoken")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestMatch(t *testing.T) {
	record := map[string]interface{}{
		"level":  "error",
		"msg":    "connection timeout occurred",
		"status": "503",
	}

	tests := []struct {
		query string
		want  bool
	}{
		{"level=error", true},
		{"level=info", false},
		{"level!=debug", true},
		{"level!=error", false},
		{"msg~timeout", true},
		{"msg~missing", false},
		{"status>200", true},
		{"status<200", false},
		{"level=error msg~timeout", true},
		{"level=error msg~missing", false},
		{"unknown=field", false},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			filters, err := Parse(tt.query)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := Match(record, filters)
			if got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}
