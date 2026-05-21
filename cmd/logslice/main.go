package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/query"
)

func main() {
	var (
		queryStr  = flag.String("query", "", "filter query (e.g. 'level == \"error\"')")
		format    = flag.String("format", "json", "output format: json, pretty, text")
		fields    = flag.String("fields", "", "comma-separated fields to include in text output")
		noColor   = flag.Bool("no-color", false, "disable color in pretty output")
	)
	flag.Parse()

	f, err := query.Parse(*queryStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: invalid query: %v\n", err)
		os.Exit(1)
	}

	var selectedFields []string
	if *fields != "" {
		selectedFields = splitFields(*fields)
	}

	cfg := output.Config{
		Format:         *format,
		SelectedFields: selectedFields,
		NoColor:        *noColor,
	}

	files := flag.Args()

	if len(files) == 0 {
		if err := process(os.Stdin, f, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for _, path := range files {
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
			os.Exit(1)
		}
		if err := process(file, f, cfg); err != nil {
			file.Close()
			fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
			os.Exit(1)
		}
		file.Close()
	}
}

func process(src *os.File, f *query.Filter, cfg output.Config) error {
	scanner := bufio.NewScanner(src)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for entry := range filter.StreamFilter(scanner, f) {
		if err := output.Write(writer, entry, cfg); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func splitFields(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if f := s[start:i]; f != "" {
				out = append(out, f)
			}
			start = i + 1
		}
	}
	return out
}
