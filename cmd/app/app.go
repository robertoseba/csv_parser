package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/robertoseba/csv_parser/internal/parser"
	"github.com/robertoseba/csv_parser/internal/printer"
	"github.com/robertoseba/csv_parser/internal/reader"
)

func Run(inputOptions *InputOptions, printer printer.IPrinter) error {
	inputReader, err := newInputReader(inputOptions.Filename)
	if err != nil {
		return err
	}
	defer inputReader.Close()

	rules, err := parser.ParseRules(inputOptions.RulesInput)
	if err != nil {
		return fmt.Errorf("error parsing rules: %w", err)
	}

	filters := parseFilters(inputOptions.FilterInput)

	config := reader.NewConfig(filters, rules)

	csvReader, err := reader.NewReader(inputReader, config)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go printer.PrintFrom(csvReader.Process(), &wg)
	wg.Wait()

	return nil
}

func parseFilters(colFilters string) []string {
	if strings.Trim(colFilters, " ") == "" {
		return nil
	}
	filters := strings.Split(colFilters, ",")
	for i, filter := range filters {
		filters[i] = strings.Trim(filter, " ")
	}
	return filters
}

func newInputReader(filename string) (*os.File, error) {
	if filename == "" {
		return readerStdin()
	}
	return readerFile(filename)
}

func readerStdin() (*os.File, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return nil, errors.New("no files or pipes provided")
	}
	return os.Stdin, nil
}

func readerFile(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}
