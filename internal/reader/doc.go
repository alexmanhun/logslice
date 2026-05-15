// Package reader provides utilities for reading and streaming structured
// JSON log records from one or more files.
//
// Each non-empty line is parsed as a JSON object. Lines that are not valid
// JSON are silently skipped, making the reader tolerant of mixed-format logs.
//
// Usage:
//
//	recs, errs := reader.ReadFiles([]string{"app.log", "worker.log"})
//	for r := range recs {
//		// process r.Fields
//	}
//	if err := <-errs; err != nil {
//		log.Fatal(err)
//	}
package reader
