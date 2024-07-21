package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

type CsvConfig struct {
	Separator    rune
	ParseNumbers bool
	ColFilters   []string
	Validator    *Validator
}

type CsvParser struct {
	config  *CsvConfig
	headers *Row
	reader  *csv.Reader
}

func New(ioReader io.Reader, config *CsvConfig) (*CsvParser, error) {
	if config == nil {
		config = &CsvConfig{
			Separator:    ',',
			ParseNumbers: false,
			ColFilters:   make([]string, 0),
			Validator:    nil,
		}
	}

	csvReader := csv.NewReader(ioReader)

	headersArr, err := csvReader.Read()

	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	headers := NewRow(headersArr, headersArr)

	if config.Validator != nil && !hasValidColumns(config.Validator.Columns(), headers) {
		return nil, errors.New("validator has invalid column")
	}

	if !hasValidColumns(config.ColFilters, headers) {
		return nil, errors.New("filter for columns has invalid column")
	}

	return &CsvParser{
		config:  config,
		reader:  csvReader,
		headers: headers,
	}, nil
}

func (r *CsvParser) ReadLine() (*Row, error) {
	var returRow *Row
	var returnError error

	for {
		recordArr, e := r.reader.Read()

		if e == io.EOF {
			returRow = nil
			returnError = io.EOF
			break
		}

		if e != nil {
			returnError = fmt.Errorf("error reading line: %w", e)
			returRow = nil
			break
		}

		row := NewRow(r.headers.Values(), recordArr)

		if r.config.Validator == nil || r.config.Validator.IsValid(row) {
			returRow = row.Only(r.config.ColFilters)
			returnError = nil
			break
		}

	}

	return returRow, returnError
}

func (r *CsvParser) Headers() *Row {
	return r.headers.Only(r.config.ColFilters)
}

func hasValidColumns(cols []string, headers *Row) bool {
	for _, col := range cols {
		if !headers.HasColumn(col) {
			return false
		}
	}
	return true
}
