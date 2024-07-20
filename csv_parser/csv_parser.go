package csv_parser

import (
	"encoding/csv"
	"fmt"
	"io"
)

type CsvParser struct {
	separator 	rune
	headers   	*Row 
	reader    	*csv.Reader
	colFilters  []string
	// validators  *Validator
}

func New(ioReader io.Reader, separator rune) (*CsvParser, error) {
	if separator == 0 {
		separator = ','
	}

	csvReader := csv.NewReader(ioReader)
	headers, err := csvReader.Read()

	// parse filters and create validator here

	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	return &CsvParser{
		separator: separator,
		reader:    csvReader,
		headers:   NewRow(headers, headers),
	}, nil
}

func (r *CsvParser) ReadLine() (*Row, error) {
	record, err := r.reader.Read()


	if err == io.EOF{
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("error reading line: %w", err)
	}

	row := NewRow(r.headers.Values(), record)
	
	return row, nil
}

func (r *CsvParser) Headers() *Row{
	return r.headers
}