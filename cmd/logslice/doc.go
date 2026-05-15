// Package main is the entry point for the logslice CLI tool.
//
// logslice streams and filters structured JSON log files using a simple
// query syntax. It reads one or more log files, applies optional filters,
// and writes matching entries to stdout in the requested format.
//
// Usage:
//
//	logslice [flags] <file...>
//
// Flags:
//
//	-query string
//	      Filter query expression (e.g. "level=error", "status!=200")
//	-format string
//	      Output format: json (default), pretty, text
//	-help
//	      Show usage information
//
// Examples:
//
//	logslice app.log
//	logslice -query 'level=error' app.log
//	logslice -format pretty -query 'status=500' access.log audit.log
package main
