package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Record represents a single parsed JSON log line with its source file.
type Record struct {
	Source string
	Line   int
	Fields map[string]interface{}
}

// ReadFile opens a file and streams parsed JSON records to the returned channel.
// The channel is closed when the file is fully read or an error occurs.
func ReadFile(path string) (<-chan Record, <-chan error) {
	records := make(chan Record, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(records)
		defer close(errs)

		f, err := os.Open(path)
		if err != nil {
			errs <- fmt.Errorf("open %s: %w", path, err)
			return
		}
		defer f.Close()

		if err := scan(path, f, records); err != nil {
			errs <- err
		}
	}()

	return records, errs
}

// ReadFiles merges records from multiple files into a single channel.
func ReadFiles(paths []string) (<-chan Record, <-chan error) {
	merged := make(chan Record, 128)
	errs := make(chan error, len(paths))

	go func() {
		defer close(merged)
		defer close(errs)

		for _, p := range paths {
			recs, ferrs := ReadFile(p)
			for r := range recs {
				merged <- r
			}
			if err := <-ferrs; err != nil {
				errs <- err
			}
		}
	}()

	return merged, errs
}

func scan(source string, r io.Reader, out chan<- Record) error {
	scanner := bufio.NewScanner(r)
	line := 0
	for scanner.Scan() {
		line++
		raw := scanner.Bytes()
		if len(raw) == 0 {
			continue
		}
		var fields map[string]interface{}
		if err := json.Unmarshal(raw, &fields); err != nil {
			// skip non-JSON lines silently
			continue
		}
		out <- Record{Source: source, Line: line, Fields: fields}
	}
	return scanner.Err()
}
