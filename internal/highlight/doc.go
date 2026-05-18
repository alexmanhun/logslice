// Package highlight provides ANSI terminal color utilities for logslice output.
//
// It maps log severity levels to colors, colorizes keys in text output, and
// provides a Strip function to remove escape codes when writing to non-terminal
// destinations such as files or pipes.
//
// Usage:
//
//	colored := highlight.ColorizeLevel("error")  // bold red "ERROR"
//	plain   := highlight.Strip(colored)           // "ERROR"
package highlight
