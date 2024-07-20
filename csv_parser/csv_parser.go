package csv_parser

import (
	"encoding/csv"
	"fmt"
	"io"
)

type CsvParser struct {
	config   	*CsvConfig
	headers   	*Row 
	reader    	*csv.Reader
	// validator 	*Validator
}

type CsvConfig struct {
	ColFilters 		[]string
	RowRules 		[]string	
	Separator 		rune
}

func New(ioReader io.Reader, config *CsvConfig) (*CsvParser, error) {
	if config == nil {
		config = &CsvConfig{
			Separator: ',',
			ColFilters: nil,
			RowRules: nil,
		}
	}
	
	if config.Separator == 0 {
		config.Separator = ','
	}

	csvReader := csv.NewReader(ioReader)
	headers, err := csvReader.Read()

	// parse filters and create validator here

	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	return &CsvParser{
		config:    config,
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

	// Parse row through validator and return filtered row
	
	return row.Only(r.config.ColFilters...), nil
}

func (r *CsvParser) FilteredHeaders() *Row{
	return r.headers.Only(r.config.ColFilters...)
}