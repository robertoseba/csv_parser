package reader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sync"

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
	readChan    chan []string
	outputChan  chan []string
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
		readChan:    make(chan []string, 100),
		outputChan:  make(chan []string, 100),
	}, nil
}
func (r *CsvReader) Process() chan []string {
	go r.read()
	r.outputChan <- r.headers.Only(r.config.ColFilters).Values()

	go r.processRecords()
	return r.outputChan
}

func (r *CsvReader) read() {
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
	wg := sync.WaitGroup{}

	for record := range r.readChan {
		wg.Add(1)
		go func() {
			row := row.NewRow(r.currentLine, r.headers.Values(), record)

			if r.config.ColRules == nil {
				r.currentLine++
				r.outputChan <- row.Only(r.config.ColFilters).Values()
			} else {
				//TODO: How do colRules interact between them? If one is valid, should we return the row?
				// Should we define the logical operator for interaction between columns? EX: (OR)col1:eq(5)||lte(10);col2:gte(10)
				// Currently we are assuming that all columns' rules must be valid to return the row

				approved := true
				for _, colRule := range r.config.ColRules {
					if !colRule.IsValid(row) {
						approved = false
					}
				}
				if approved {
					r.currentLine++
					r.outputChan <- row.Only(r.config.ColFilters).Values()
				}
			}
			wg.Done()
		}()
	}
	close(r.outputChan)
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
