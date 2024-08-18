package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/robertoseba/csv_parser/internal/parser"
	"github.com/robertoseba/csv_parser/internal/printer"
	"github.com/robertoseba/csv_parser/internal/reader"
)

func Run(ioReader io.Reader, colFilters string, rowRules string) {
	rules, err := parser.ParseRules(rowRules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing rules: %s\n", err)
		os.Exit(1)
	}

	csvConfig := &reader.CsvConfig{
		ColFilters: splitFilters(colFilters),
		ColRules:   rules,
	}

	csvReader, err := reader.NewParser(ioReader, csvConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating csv: %s\n", err)
		os.Exit(1)
	}

	printer := printer.NewPrinter()

	printer.PrintHeader(csvReader.Headers().Values())

	for {
		row, err := csvReader.ReadLine()

		if errors.Is(err, io.EOF) {
			break
		}

		if errors.Is(err, reader.ErrInvalidRow) {
			continue
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unexpected error reading line: %s\n", err)
			os.Exit(1)
		}
		printer.PrintRow(row.Values())
	}
}

func splitFilters(colFilters string) []string {
	if strings.Trim(colFilters, " ") == "" {
		return nil
	}

	filters := strings.Split(colFilters, ",")
	for i, filter := range filters {
		filters[i] = strings.Trim(filter, " ")
	}

	return filters
}
