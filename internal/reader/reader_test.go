package reader

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/robertoseba/csv_parser/internal/parser"
)

func TestReaderHeaders(t *testing.T) {

	config := &CsvConfig{}

	testReader := strings.NewReader("col1,col2,col3\nrow_1000,2,3\nrow_11,5,6\nrow_99,8,9")

	expected := []string{"col1", "col2", "col3"}

	results, err := NewReader(testReader, config)

	if err != nil {
		t.Errorf("Failted creating parser: %v", err)
	}

	if !slices.Equal(results.headers.Values(), expected) {
		t.Errorf("Expected %v, got %v", expected, results.headers.Values())
	}
}

func TestReaderColFilters(t *testing.T) {
	tests := []struct {
		name        string
		inputConfig *CsvConfig
		expected    []string
		err         error
	}{
		{
			name:        "ColFilters col1 and col2",
			inputConfig: &CsvConfig{ColFilters: []string{"col1", "col2"}},
			expected:    []string{"col1", "col2"},
			err:         nil,
		},
		{
			name:        "No colFilters",
			inputConfig: &CsvConfig{},
			expected:    nil,
			err:         nil,
		},
		{
			name: "Error when ColFilters not present in column headers",
			inputConfig: &CsvConfig{
				ColFilters: []string{"col4", "col5"},
			},
			expected: nil,
			err:      errors.New("filter for columns has invalid column"),
		},
	}

	for _, test := range tests {
		testReader := strings.NewReader("col1,col2,col3\nrow_1000,2,3\nrow_11,5,6\nrow_99,8,9")

		t.Run(test.name, func(t *testing.T) {

			results, err := NewReader(testReader, test.inputConfig)

			if err != nil {
				if err.Error() != test.err.Error() {
					t.Errorf("Failed creating parser: %v", err)
				}
			}

			if results != nil {
				if !slices.Equal(results.config.ColFilters, test.expected) {
					t.Errorf("Expected %v, got %v", test.expected, results.config.ColFilters)
				}
			}
		})
	}
}

func TestReaderReadLine(t *testing.T) {
	testRule, _ := parser.ParseRules("col2:gte(8)")

	tests := []struct {
		name        string
		inputConfig *CsvConfig
		expected    [][]string
		errs        []error
	}{
		{
			name:        "reads all lines with no filter nor rules",
			inputConfig: &CsvConfig{},
			expected: [][]string{
				{"col1", "col2", "col3"},
				{"row_1000", "2", "3"},
				{"row_11", "5", "6"},
				{"row_99", "8", "9"},
			},
		},
		{
			name:        "readlines with col1 and col2 filters",
			inputConfig: &CsvConfig{ColFilters: []string{"col1", "col2"}},
			expected: [][]string{
				{"col1", "col2"},
				{"row_1000", "2"},
				{"row_11", "5"},
				{"row_99", "8"},
			},
		},
		{
			name: "returns error when line fails rules, returning only valid rows",
			inputConfig: &CsvConfig{
				ColRules: testRule,
			},
			expected: [][]string{
				{"col1", "col2", "col3"},
				{"row_99", "8", "9"},
			},
		},
	}

	for _, test := range tests {
		testReader := strings.NewReader("col1,col2,col3\nrow_1000,2,3\nrow_11,5,6\nrow_99,8,9")

		t.Run(test.name, func(t *testing.T) {

			reader, err := NewReader(testReader, test.inputConfig)

			if err != nil {
				t.Errorf("Failed creating parser: %v", err)
			}

			i := 0
			for res := range reader.Process() {
				if res != nil && !slices.Equal(res, test.expected[i]) {
					t.Errorf("Expected %v, got %v", test.expected[i], res)
				}
				i++
			}
		})
	}
}
