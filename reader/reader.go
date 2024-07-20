package reader

import (
	"encoding/csv"
	"errors"
	"os"
)

type Reader struct {
	filePath  string
	separator rune
	headers   []string
	reader    *csv.Reader
}

func New(filePath string, separator rune) (*Reader, error) {
	if separator == 0 {
		separator = ','
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("error opening file")
	}

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	// parse filters and create validator here
	if err != nil {
		return nil, errors.New("error reading headers")
	}

	return &Reader{
		filePath:  filePath,
		separator: separator,
		reader:    reader,
		headers:   headers,
	}, nil
}

func (r *Reader) ReadLine() ([]string, error) {
	return r.reader.Read()
}
