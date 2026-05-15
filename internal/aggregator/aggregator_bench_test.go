package aggregator

import (
	"fmt"
	"testing"
)

func benchmarkEntries(n int) []map[string]interface{} {
	levels := []string{"info", "warn", "error", "debug"}
	entries := make([]map[string]interface{}, n)
	for i := 0; i < n; i++ {
		entries[i] = map[string]interface{}{
			"level":   levels[i%len(levels)],
			"service": fmt.Sprintf("svc-%d", i%10),
			"msg":     fmt.Sprintf("log message number %d", i),
		}
	}
	return entries
}

func BenchmarkCountBy(b *testing.B) {
	entries := benchmarkEntries(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountBy(entries, "level")
	}
}

func BenchmarkUniqueValues(b *testing.B) {
	entries := benchmarkEntries(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UniqueValues(entries, "service")
	}
}
