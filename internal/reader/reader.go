package reader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/robertoseba/csv_parser/internal/parser"
	"github.com/robertoseba/csv_parser/internal/row"
)

var ErrInvalidRow = errors.New("invalid row")

type CsvConfig struct {
	ColFilters []string
	ColRules   []parser.ColRules
}

type CsvReader struct {
	currentLine int
	config      *CsvConfig
	headers     *row.Row
	reader      *csv.Reader
}

func NewReader(ioReader io.Reader, config *CsvConfig) (*CsvReader, error) {
	if config == nil {
		config = &CsvConfig{
			ColFilters: make([]string, 0),
			ColRules:   nil,
		}
	}

	csvReader := csv.NewReader(ioReader)

	headersArr, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	headers := row.NewRow(0, headersArr, headersArr)

	if config.ColRules != nil && !isColRulesValid(config.ColRules, headers) {
		return nil, errors.New("rules have invalid column")
	}

	if !isFilterColsValid(config.ColFilters, headers) {
		return nil, errors.New("filter for columns has invalid column")
	}

	return &CsvReader{
		currentLine: 1,
		config:      config,
		reader:      csvReader,
		headers:     headers,
	}, nil
}

func (r *CsvReader) ReadLine() (*row.Row, error) {
	recordArr, err := r.reader.Read()

	if errors.Is(err, io.EOF) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("unexpected error reading line: %w", err)
	}

	row := row.NewRow(r.currentLine, r.headers.Values(), recordArr)

	if r.config.ColRules == nil {
		r.currentLine++
		return row.Only(r.config.ColFilters), nil
	}

	//TODO: How do colRules interact between them? If one is valid, should we return the row?
	// Should we define the logical operator for interaction between columns? EX: (OR)col1:eq(5)||lte(10);col2:gte(10)
	// Currently we are assuming that all columns' rules must be valid to return the row
	for _, colRule := range r.config.ColRules {
		if !colRule.IsValid(row) {
			return nil, ErrInvalidRow
		}
	}

	r.currentLine++
	return row.Only(r.config.ColFilters), nil
}

func (r *CsvReader) Headers() *row.Row {
	return r.headers.Only(r.config.ColFilters)
}

func isColRulesValid(colRules []parser.ColRules, headers *row.Row) bool {
	for _, colRule := range colRules {
		if !headers.HasColumn(colRule.Column()) {
			return false
		}
	}
	return true
}

func isFilterColsValid(cols []string, headers *row.Row) bool {
	for _, col := range cols {
		if !headers.HasColumn(col) {
			return false
		}
	}
	return true
}
