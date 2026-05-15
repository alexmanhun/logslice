// Package sorter provides sorting utilities for structured log entries.
//
// It operates on slices of map[string]any, which is the native representation
// used when unmarshalling JSON log lines. Entries can be sorted ascending or
// descending by any top-level string or numeric field.
//
// Example usage:
//
//	err := sorter.SortEntries(entries, "ts", sorter.Ascending)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Entries that do not contain the specified field are placed at the end of the
// result regardless of the chosen sort order.
package sorter
