package limiter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/limiter"
)

func TestLimiter_Allow_WithLimit(t *testing.T) {
	l := limiter.New(3)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() to return true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow() to return false after limit reached")
	}
}

func TestLimiter_Allow_Unlimited(t *testing.T) {
	l := limiter.New(0)
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatalf("expected unlimited limiter to always allow, failed at %d", i)
		}
	}
}

func TestLimiter_Remaining(t *testing.T) {
	l := limiter.New(5)
	if l.Remaining() != 5 {
		t.Fatalf("expected 5 remaining, got %d", l.Remaining())
	}
	l.Allow()
	l.Allow()
	if l.Remaining() != 3 {
		t.Fatalf("expected 3 remaining after 2 allows, got %d", l.Remaining())
	}
}

func TestLimiter_Remaining_Unlimited(t *testing.T) {
	l := limiter.New(-1)
	if l.Remaining() != -1 {
		t.Fatalf("expected -1 for unlimited, got %d", l.Remaining())
	}
}

func TestLimiter_Reset(t *testing.T) {
	l := limiter.New(2)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected false after limit exhausted")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected true after reset")
	}
}

func TestApply_LimitsEntries(t *testing.T) {
	entries := []map[string]any{
		{"id": 1}, {"id": 2}, {"id": 3}, {"id": 4}, {"id": 5},
	}
	result := limiter.Apply(entries, 3)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestApply_NoLimit(t *testing.T) {
	entries := []map[string]any{
		{"id": 1}, {"id": 2}, {"id": 3},
	}
	result := limiter.Apply(entries, 0)
	if len(result) != 3 {
		t.Fatalf("expected all 3 entries, got %d", len(result))
	}
}

func TestApply_LimitExceedsEntries(t *testing.T) {
	entries := []map[string]any{
		{"id": 1}, {"id": 2},
	}
	result := limiter.Apply(entries, 10)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}
