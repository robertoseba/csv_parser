package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/robertoseba/csv_parser/pkg/row"
	"github.com/robertoseba/csv_parser/pkg/rule"
)

type CsvConfig struct {
	Separator    rune
	ParseNumbers bool
	ColFilters   []string
	ColRules     []*rule.ColRules
}

type CsvParser struct {
	config  *CsvConfig
	headers *row.Row
	reader  *csv.Reader
}

func NewParser(ioReader io.Reader, config *CsvConfig) (*CsvParser, error) {
	if config == nil {
		config = &CsvConfig{
			Separator:  ',',
			ColFilters: make([]string, 0),
			ColRules:   nil,
		}
	}

	csvReader := csv.NewReader(ioReader)

	headersArr, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	headers := row.NewRow(headersArr, headersArr)

	if config.ColRules != nil && !isColRulesValid(config.ColRules, headers) {
		return nil, errors.New("validator has invalid column")
	}

	if !isFilterColsValid(config.ColFilters, headers) {
		return nil, errors.New("filter for columns has invalid column")
	}

	return &CsvParser{
		config:  config,
		reader:  csvReader,
		headers: headers,
	}, nil
}

func (r *CsvParser) ReadLine() (*row.Row, error) {
	recordArr, e := r.reader.Read()

	if e == io.EOF {
		return nil, io.EOF
	}

	if e != nil {
		return nil, fmt.Errorf("unexpected error reading line: %w", e)
	}

	row := row.NewRow(r.headers.Values(), recordArr)

	if r.config.ColRules == nil {
		return row.Only(r.config.ColFilters), nil
	}

	//TODO: How do colRules interact between them? If one is valid, should we return the row?
	// Should we define the logical operator for interaction between columns? EX: (OR)col1:eq(5)ANDlte(10);col2:gte(10)
	for _, colRule := range r.config.ColRules {
		if colRule.IsValid(row) {
			return row.Only(r.config.ColFilters), nil
		}
	}

	return nil, ErrInvalidRow
}

func (r *CsvParser) Headers() *row.Row {
	return r.headers.Only(r.config.ColFilters)
}

func isColRulesValid(colRules []*rule.ColRules, headers *row.Row) bool {
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
