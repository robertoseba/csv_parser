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
	HeaderOnly bool
}

func NewConfig(colFilters []string, colRules []parser.ColRules, headerOnly bool) *CsvConfig {
	return &CsvConfig{
		ColFilters: colFilters,
		ColRules:   colRules,
		HeaderOnly: headerOnly,
	}
}

type CsvReader struct {
	currentLine int
	config      *CsvConfig
	headers     *row.Row
	reader      *csv.Reader
	readChan    chan []string
	outputChan  chan []string
}

func NewReader(ioReader io.Reader, config *CsvConfig) (*CsvReader, error) {
	if config == nil {
		config = &CsvConfig{
			ColFilters: nil,
			ColRules:   nil,
			HeaderOnly: false,
		}
	}

	csvReader := csv.NewReader(ioReader)

	headersArr, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	headers := row.NewRow(0, headersArr, headersArr)

	if config.ColRules != nil && !isColRulesValid(config.ColRules, headers) {
		return nil, errors.New("rules have one or more invalid columns")
	}

	if !isFilterColsValid(config.ColFilters, headers) {
		return nil, errors.New("filter for columns has invalid column")
	}

	return &CsvReader{
		currentLine: 1,
		config:      config,
		reader:      csvReader,
		headers:     headers,
		readChan:    make(chan []string, 150),
		outputChan:  make(chan []string, 150),
	}, nil
}

func (r *CsvReader) Process() chan []string {
	if r.config.HeaderOnly {
		r.outputChan <- []string{"Headers"}
		for _, header := range r.headers.Values() {
			r.outputChan <- []string{header}
		}
		close(r.outputChan)
		return r.outputChan
	}

	r.outputChan <- r.headers.Only(r.config.ColFilters).Values()

	go r.getDataFromReader()

	go r.processRecords()

	return r.outputChan
}

func (r *CsvReader) getDataFromReader() {
	for {
		record, err := r.reader.Read()
		if errors.Is(err, io.EOF) {
			close(r.readChan)
			break
		}

		if err != nil {
			close(r.readChan)
			panic(err)
		}
		r.readChan <- record
	}
}

func (r *CsvReader) processRecords() {
	defer close(r.outputChan)

	for record := range r.readChan {
		row := row.NewRow(r.currentLine, r.headers.Values(), record)

		if r.config.ColRules == nil {
			r.currentLine++
			r.outputChan <- row.Only(r.config.ColFilters).Values()
		} else {
			// Should we define the logical operator for interaction between columns? EX: (OR)col1:eq(5)||lte(10);col2:gte(10)
			// Currently we are assuming that all columns' rules must be valid to return the row
			isValid := true
			for _, colRule := range r.config.ColRules {
				if !colRule.IsValid(row) {
					isValid = false
					break
				}
			}
			if isValid {
				r.currentLine++
				r.outputChan <- row.Only(r.config.ColFilters).Values()
			}
		}
	}
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
