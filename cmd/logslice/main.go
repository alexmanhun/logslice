package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/query"
	"github.com/yourorg/logslice/internal/reader"
)

func main() {
	var (
		queryStr  = flag.String("query", "", "Filter query (e.g. 'level=error')")
		format    = flag.String("format", "json", "Output format: json, pretty, text")
		showHelp  = flag.Bool("help", false, "Show usage")
	)
	flag.Parse()

	if *showHelp || flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: logslice [flags] <file...>\n\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	var filter *query.Filter
	if *queryStr != "" {
		var err error
		filter, err = query.Parse(*queryStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logslice: invalid query: %v\n", err)
			os.Exit(1)
		}
	}

	formatter, err := output.New(*format, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}

	entries, errs := reader.ReadFiles(flag.Args())

	for _, entry := range entries {
		if filter == nil || query.Match(filter, entry) {
			if writeErr := formatter.Write(entry); writeErr != nil {
				fmt.Fprintf(os.Stderr, "logslice: write error: %v\n", writeErr)
			}
		}
	}

	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", e)
	}

	if len(errs) > 0 {
		os.Exit(2)
	}
}
