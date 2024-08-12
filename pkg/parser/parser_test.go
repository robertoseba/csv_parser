package parser

import (
	"encoding/csv"
	"errors"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/robertoseba/csv_parser/pkg/row"
	"github.com/robertoseba/csv_parser/pkg/rule"
)

func TestParserHeaders(t *testing.T) {
	colRules, _ := rule.NewFrom("col1:eq(row_1000)||eq(row_11)")

	config := &CsvConfig{
		ColFilters: []string{"col1", "col2"},
		ColRules:   colRules,
	}

	testReader := strings.NewReader("col1,col2,col3\nrow_1000,2,3\nrow_11,5,6\nrow_99,8,9")

	expected := &CsvParser{
		config:  config,
		reader:  csv.NewReader(testReader),
		headers: row.NewRow([]string{"col1", "col2", "col3"}, []string{"col1", "col2", "col3"}),
	}

	results, err := NewParser(testReader, config)

	if err != nil {
		t.Errorf("Failted creating parser: %v", err)
	}

	if !slices.Equal(results.Headers().Values(), expected.Headers().Values()) {
		t.Errorf("Expected %v, got %v", expected.Headers(), results.Headers())
	}
}

func TestParserColFilters(t *testing.T) {
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
			name: "ColFilters not present in headers",
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

			results, err := NewParser(testReader, test.inputConfig)

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

func TestParserReadLine(t *testing.T) {
	testRule, _ := rule.NewFrom("col2:gte(8)")

	tests := []struct {
		name        string
		inputConfig *CsvConfig
		expected    [][]string
		errs        []error
	}{
		{
			name:        "readlines with col1 and col2 filters",
			inputConfig: &CsvConfig{ColFilters: []string{"col1", "col2"}},
			expected: [][]string{
				{"row_1000", "2"},
				{"row_11", "5"},
				{"row_99", "8"},
			},
			errs: nil,
		},
		{
			name:        "reads all lines with no filter nor rules",
			inputConfig: &CsvConfig{},
			expected: [][]string{
				{"row_1000", "2", "3"},
				{"row_11", "5", "6"},
				{"row_99", "8", "9"},
			},
			errs: nil,
		},
		{
			name: "returns error when line fails rules, returning only valid rows",
			inputConfig: &CsvConfig{
				ColRules: testRule,
			},
			expected: [][]string{
				nil,
				nil,
				{"row_99", "8", "9"},
			},
			errs: []error{
				ErrInvalidRow,
				ErrInvalidRow,
			},
		},
	}

	for _, test := range tests {
		testReader := strings.NewReader("col1,col2,col3\nrow_1000,2,3\nrow_11,5,6\nrow_99,8,9")

		t.Run(test.name, func(t *testing.T) {

			results, err := NewParser(testReader, test.inputConfig)

			if err != nil {
				t.Errorf("Failed creating parser: %v", err)
			}

			i := 0

			for {
				row, err := results.ReadLine()

				if errors.Is(err, io.EOF) {
					break
				}

				if err != nil {
					if !errors.Is(err, test.errs[i]) {
						t.Errorf("Expected %v, got %v", test.errs[i], err)
					}
				}

				if row != nil && !slices.Equal(row.Values(), test.expected[i]) {
					t.Errorf("Expected %v, got %v", test.expected[i], row.Values())
				}

				i += 1
			}
		})
	}
}
