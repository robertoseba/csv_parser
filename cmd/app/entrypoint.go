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

func Run(filePath string, colFilters string, rowRules string) {
	file := openFile(filePath)
	defer file.Close()

	rules, err := rule.NewFrom(rowRules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing rules: %s\n", err)
		os.Exit(1)
	}

	csvConfig := &parser.CsvConfig{
		Separator:  ',',
		ColFilters: splitStringColFilters(colFilters),
		ColRules:   rules,
	}

	reader, err := parser.NewParser(file, csvConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating csv: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(reader.Headers())

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

		fmt.Println(row)
	}
}

func splitStringColFilters(colFilters string) []string {
	if strings.Trim(colFilters, " ") == "" {
		return nil
	}
	return strings.Split(colFilters, ",")
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", filePath)
		os.Exit(1)
	}
	return file
}
