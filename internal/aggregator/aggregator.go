package aggregator

import (
	"fmt"
	"sort"
)

// CountResult holds the count for a specific field value.
type CountResult struct {
	Value string
	Count int
}

// CountBy counts log entries grouped by the given field name.
func CountBy(entries []map[string]interface{}, field string) []CountResult {
	counts := make(map[string]int)

	for _, entry := range entries {
		val, ok := entry[field]
		if !ok {
			counts["<missing>"]++
			continue
		}
		counts[fmt.Sprintf("%v", val)]++
	}

	results := make([]CountResult, 0, len(counts))
	for v, c := range counts {
		results = append(results, CountResult{Value: v, Count: c})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Value < results[j].Value
	})

	return results
}

// UniqueValues returns the distinct values for the given field across entries.
func UniqueValues(entries []map[string]interface{}, field string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, entry := range entries {
		val, ok := entry[field]
		if !ok {
			continue
		}
		s := fmt.Sprintf("%v", val)
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}

	sort.Strings(result)
	return result
}
