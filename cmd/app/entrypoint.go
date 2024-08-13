package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/robertoseba/csv_parser/pkg/parser"
	"github.com/robertoseba/csv_parser/pkg/rule"
)

func Run(input io.Reader, colFilters string, rowRules string) {
	rules, err := rule.NewFrom(rowRules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing rules: %s\n", err)
		os.Exit(1)
	}

	csvConfig := &parser.CsvConfig{
		ColFilters: splitStringColFilters(colFilters),
		ColRules:   rules,
	}

	reader, err := parser.NewParser(input, csvConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating csv: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf(" \t%s\n", reader.Headers())

	for {
		row, err := reader.ReadLine()

		if errors.Is(err, io.EOF) {
			break
		}

		if errors.Is(err, parser.ErrInvalidRow) {
			continue
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unexpected error reading line: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%d\t%s\n", row.LineNumber(), row)
	}
}

func splitStringColFilters(colFilters string) []string {
	if strings.Trim(colFilters, " ") == "" {
		return nil
	}

	filters := strings.Split(colFilters, ",")
	for i, filter := range filters {
		filters[i] = strings.Trim(filter, " ")
	}

	return filters
}
