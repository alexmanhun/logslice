package query

import (
	"fmt"
	"strconv"
	"strings"
)

// Op represents a comparison operator in a filter expression.
type Op string

const (
	OpEq      Op = "="
	OpNeq     Op = "!="
	OpContains Op = "~"
	OpGt      Op = ">"
	OpLt      Op = "<"
)

// Filter represents a single parsed filter expression, e.g. level=error
type Filter struct {
	Field string
	Op    Op
	Value string
}

// Parse parses a query string into a slice of Filters.
// Query syntax: "field=value field2~substring field3!=foo"
func Parse(query string) ([]Filter, error) {
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}

	tokens := strings.Fields(query)
	filters := make([]Filter, 0, len(tokens))

	for _, token := range tokens {
		f, err := parseToken(token)
		if err != nil {
			return nil, fmt.Errorf("invalid filter %q: %w", token, err)
		}
		filters = append(filters, f)
	}

	return filters, nil
}

func parseToken(token string) (Filter, error) {
	for _, op := range []Op{OpNeq, OpContains, OpGt, OpLt, OpEq} {
		if idx := strings.Index(token, string(op)); idx > 0 {
			return Filter{
				Field: token[:idx],
				Op:    op,
				Value: token[idx+len(op):],
			}, nil
		}
	}
	return Filter{}, fmt.Errorf("no valid operator found")
}

// Match reports whether the given log record (as a map) satisfies all filters.
func Match(record map[string]interface{}, filters []Filter) bool {
	for _, f := range filters {
		val, ok := record[f.Field]
		if !ok {
			return false
		}
		strVal := fmt.Sprintf("%v", val)
		if !matchOp(strVal, f.Op, f.Value) {
			return false
		}
	}
	return true
}

func matchOp(fieldVal string, op Op, filterVal string) bool {
	switch op {
	case OpEq:
		return fieldVal == filterVal
	case OpNeq:
		return fieldVal != filterVal
	case OpContains:
		return strings.Contains(fieldVal, filterVal)
	case OpGt:
		a, err1 := strconv.ParseFloat(fieldVal, 64)
		b, err2 := strconv.ParseFloat(filterVal, 64)
		if err1 != nil || err2 != nil {
			return fieldVal > filterVal
		}
		return a > b
	case OpLt:
		a, err1 := strconv.ParseFloat(fieldVal, 64)
		b, err2 := strconv.ParseFloat(filterVal, 64)
		if err1 != nil || err2 != nil {
			return fieldVal < filterVal
		}
		return a < b
	}
	return false
}
