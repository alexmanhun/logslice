// Package output provides log entry formatting for logslice.
//
// Supported formats:
//
//   - json   — compact single-line JSON (default)
//   - pretty — indented multi-line JSON
//   - text   — key=value pairs sorted alphabetically
//
// Example usage:
//
//	f := output.New(os.Stdout, output.FormatText)
//	f.Write(map[string]any{"level": "info", "msg": "started"})
package output
