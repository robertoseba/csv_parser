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
	ColFilters  []string
	ColRules    []parser.ColRules
	OrderedRows bool
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
			ColFilters:  make([]string, 0),
			ColRules:    nil,
			OrderedRows: false,
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
		readChan:    make(chan []string, 150),
		outputChan:  make(chan []string, 150),
	}, nil
}
func (r *CsvReader) Process() chan []string {
	go r.read()
	r.outputChan <- r.headers.Only(r.config.ColFilters).Values()

	numWorkers := 10
	if r.config.OrderedRows {
		numWorkers = 1
	}

	wg := sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go r.processRecords(&wg)
	}
	go r.done(&wg)

	return r.outputChan
}

func (r *CsvReader) done(wg *sync.WaitGroup) {
	defer close(r.outputChan)
	wg.Wait()
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

func (r *CsvReader) processRecords(wg *sync.WaitGroup) {
	for record := range r.readChan {
		row := row.NewRow(r.currentLine, r.headers.Values(), record)

		if r.config.ColRules == nil {
			r.currentLine++
			r.outputChan <- row.Only(r.config.ColFilters).Values()
		} else {
			//TODO: How do colRules interact between them? If one is valid, should we return the row?
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
	wg.Done()
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
