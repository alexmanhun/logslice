// Package limiter provides utilities for limiting the number of log entries
// returned during streaming and filtering operations.
package limiter

// Limiter controls how many entries are passed through a processing pipeline.
type Limiter struct {
	max   int
	count int
}

// New creates a new Limiter with the given maximum number of entries.
// If max is <= 0, the limiter is considered unlimited.
func New(max int) *Limiter {
	return &Limiter{max: max}
}

// Allow reports whether the next entry should be allowed through.
// Once the limit is reached, Allow returns false for all subsequent calls.
func (l *Limiter) Allow() bool {
	if l.max <= 0 {
		return true
	}
	if l.count >= l.max {
		return false
	}
	l.count++
	return true
}

// Remaining returns the number of entries still allowed.
// Returns -1 if the limiter is unlimited.
func (l *Limiter) Remaining() int {
	if l.max <= 0 {
		return -1
	}
	remaining := l.max - l.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset resets the limiter's internal counter back to zero.
func (l *Limiter) Reset() {
	l.count = 0
}

// Apply filters entries using the limiter, returning only those allowed.
func Apply(entries []map[string]any, max int) []map[string]any {
	if max <= 0 {
		return entries
	}
	l := New(max)
	result := make([]map[string]any, 0, max)
	for _, entry := range entries {
		if !l.Allow() {
			break
		}
		result = append(result, entry)
	}
	return result
}
