// Package aggregator provides utilities for summarising structured log entries.
//
// It supports counting entries grouped by a specific field and extracting
// unique field values across a set of log entries.
//
// Example usage:
//
//	results := aggregator.CountBy(entries, "level")
//	for _, r := range results {
//		fmt.Printf("%s: %d\n", r.Value, r.Count)
//	}
//
//	uniq := aggregator.UniqueValues(entries, "service")
package aggregator
