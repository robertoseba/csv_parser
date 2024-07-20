package reader

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Reader struct {
	separator rune
	headers   *Row 
	reader    *csv.Reader
}

func New(ioReader io.Reader, separator rune) (*Reader, error) {
	if separator == 0 {
		separator = ','
	}

	csvReader := csv.NewReader(ioReader)
	headers, err := csvReader.Read()

	// parse filters and create validator here

	if err != nil {
		return nil, fmt.Errorf("error parsing headers: %w", err)
	}

	return &Reader{
		separator: separator,
		reader:    csvReader,
		headers:   NewRow(headers, headers),
	}, nil
}

func (r *Reader) ReadLine() (*Row, error) {
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

func (r *Reader) Headers() *Row{
	return r.headers
}