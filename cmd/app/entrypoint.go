package app

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/robertoseba/csv_parser/internal/parser"
	"github.com/robertoseba/csv_parser/internal/printer"
	"github.com/robertoseba/csv_parser/internal/reader"
)

func Run(filename string, colFilters string, rowRules string) {
	var err error

	rules, err := parser.ParseRules(rowRules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing rules: %s\n", err)
		os.Exit(1)
	}

	csvConfig := &reader.CsvConfig{
		ColFilters: splitFilters(colFilters),
		ColRules:   rules,
	}

	var csvReader *reader.CsvReader

	if filename == "" {
		csvReader, err = reader.NewReader(readerStdin(), csvConfig)
	} else {
		f := readerFile(filename)
		csvReader, err = reader.NewReader(f, csvConfig)
		defer f.Close()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating csv: %s\n", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	printer := printer.NewPrinter(csvReader.Process(), &wg)
	go printer.Start()

	wg.Wait()
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

func readerStdin() *os.File {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("No files or pipes provided")
		os.Exit(1)
	}
	return os.Stdin
}

func readerFile(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file: err:", err)
		os.Exit(1)
	}
	return f
}
