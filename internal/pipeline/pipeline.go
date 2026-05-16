// Package pipeline wires together reading, filtering, sorting, limiting,
// and formatting into a single reusable processing chain.
package pipeline

import (
	"io"

	"github.com/yourorg/logslice/internal/aggregator"
	"github.com/yourorg/logslice/internal/limiter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/query"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/sorter"
)

// Config holds all options that control a pipeline run.
type Config struct {
	Files     []string
	Filter    string
	SortField string
	SortDesc  bool
	Limit     int
	Format    string // "json", "pretty", "text"
	CountBy   string // when non-empty, emit aggregation instead of entries
}

// Run executes the full pipeline and writes results to w.
func Run(cfg Config, w io.Writer) error {
	// 1. Parse filter (empty filter is valid — matches everything).
	filter, err := query.Parse(cfg.Filter)
	if err != nil {
		return err
	}

	// 2. Read entries from all files.
	entries, err := reader.ReadFiles(cfg.Files)
	if err != nil {
		return err
	}

	// 3. Filter.
	var matched []map[string]any
	for _, e := range entries {
		if filter == nil || query.Match(e, filter) {
			matched = append(matched, e)
		}
	}

	// 4. Sort (optional).
	if cfg.SortField != "" {
		sorter.SortEntries(matched, cfg.SortField, cfg.SortDesc)
	}

	// 5. Limit.
	lim := limiter.New(cfg.Limit)
	matched = lim.Apply(matched)

	// 6. Aggregation shortcut — count-by replaces normal output.
	if cfg.CountBy != "" {
		counts := aggregator.CountBy(matched, cfg.CountBy)
		fmt := output.New(w, "json")
		for k, v := range counts {
			_ = fmt.Write(map[string]any{cfg.CountBy: k, "count": v})
		}
		return nil
	}

	// 7. Format and write.
	fmt := output.New(w, cfg.Format)
	for _, e := range matched {
		if err := fmt.Write(e); err != nil {
			return err
		}
	}
	return nil
}
