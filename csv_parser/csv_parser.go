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
	validator 	*Validator
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
	
	headersArr, err := csvReader.Read()

	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	headers := NewRow(headersArr, headersArr)

	// TODO: return which col filter is wrong
	if !isColFiltersValid(config.ColFilters, headers){
		return nil, fmt.Errorf("error parsing Col filters")
	}

	validator, err := NewValidator(config.RowRules, headers)
	
	if err != nil {
		return nil, fmt.Errorf("error creating validator: %w", err)
	}

	return &CsvParser{
		config:    config,
		reader:    csvReader,
		headers:   headers,
		validator: validator,
	}, nil
}

func (r *CsvParser) ReadLine() (*Row, error) {
	var returRow *Row 
	var returnError error

	for {
		recordArr, e := r.reader.Read()

		if e == io.EOF{
			returRow = nil
			returnError = io.EOF
			break	
		}

		if e != nil {
			returnError = fmt.Errorf("error reading line: %w", e)
			returRow = nil
			break
		}

		row:= NewRow(r.headers.Values(), recordArr)

		if r.validator.IsValid(row){
			returRow = row.Only(r.config.ColFilters...)
			returnError = nil
			break
		}
		
	}

	return returRow, returnError 
}

func (r *CsvParser) FilteredHeaders() *Row{
	return r.headers.Only(r.config.ColFilters...)
}

func isColFiltersValid(colFilters []string, headers *Row) bool{
	result := true 

	for _, col := range colFilters {
		if !headers.Contains(col){
			result = false
			break
		}
	}
	return result 
}