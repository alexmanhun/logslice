// Package highlight provides terminal color highlighting for log output.
package highlight

import (
	"fmt"
	"strings"
)

// ANSI color codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
	Gray   = "\033[90m"
)

// LevelColor returns an ANSI color code for the given log level string.
func LevelColor(level string) string {
	switch strings.ToLower(level) {
	case "error", "fatal", "critical":
		return Red
	case "warn", "warning":
		return Yellow
	case "info":
		return Green
	case "debug", "trace":
		return Cyan
	default:
		return Gray
	}
}

// ColorizeLevel wraps the level string in its corresponding ANSI color.
func ColorizeLevel(level string) string {
	color := LevelColor(level)
	return fmt.Sprintf("%s%s%s", color, strings.ToUpper(level), Reset)
}

// ColorizeKey wraps a key name in bold for display.
func ColorizeKey(key string) string {
	return fmt.Sprintf("%s%s%s", Bold, key, Reset)
}

// Strip removes all ANSI escape codes from a string.
func Strip(s string) string {
	var b strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
