// Package pipeline provides a high-level Run function that composes the
// reader, query, sorter, limiter, aggregator, and output packages into a
// single processing chain.
//
// Typical usage:
//
//	err := pipeline.Run(pipeline.Config{
//		Files:  []string{"app.log"},
//		Filter: "level=error",
//		Limit:  50,
//		Format: "pretty",
//	}, os.Stdout)
//
// Config fields are all optional; a zero-value Config reads from no files
// and writes nothing, which is a no-op.
package pipeline
