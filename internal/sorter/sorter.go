// Package sorter provides utilities for sorting log entries by a specified field.
package sorter

import (
	"fmt"
	"sort"
)

// Order defines the sort direction.
type Order int

const (
	Ascending Order = iota
	Descending
)

// SortEntries sorts a slice of log entry maps by the given field.
// Entries missing the field are placed at the end.
func SortEntries(entries []map[string]any, field string, order Order) error {
	if field == "" {
		return fmt.Errorf("sort field must not be empty")
	}

	var sortErr error

	sort.SliceStable(entries, func(i, j int) bool {
		if sortErr != nil {
			return false
		}

		vi, oki := entries[i][field]
		vj, okj := entries[j][field]

		// Entries missing the field sink to the bottom.
		if !oki && !okj {
			return false
		}
		if !oki {
			return false
		}
		if !okj {
			return true
		}

		cmp, err := compareValues(vi, vj)
		if err != nil {
			sortErr = fmt.Errorf("field %q: %w", field, err)
			return false
		}

		if order == Descending {
			return cmp > 0
		}
		return cmp < 0
	})

	return sortErr
}

// compareValues compares two values of comparable types.
// Supports string, float64 (JSON numbers), and bool.
func compareValues(a, b any) (int, error) {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		if !ok {
			return 0, fmt.Errorf("type mismatch: string vs %T", b)
		}
		if av < bv {
			return -1, nil
		}
		if av > bv {
			return 1, nil
		}
		return 0, nil
	case float64:
		bv, ok := b.(float64)
		if !ok {
			return 0, fmt.Errorf("type mismatch: float64 vs %T", b)
		}
		if av < bv {
			return -1, nil
		}
		if av > bv {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type %T for comparison", a)
	}
}
