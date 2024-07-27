package app

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/robertoseba/csv_parser/pkg/parser"
)

func Run(filePath string, colFilters string, rowRules string) {
	file := openFile(filePath)
	defer file.Close()

	rules := createRules(rowRules)
	var validator *parser.Validator
	if rules != nil {
		validator = parser.NewValidator(rules)
	}

	csvConfig := &parser.CsvConfig{
		Separator:  ',',
		ColFilters: splitStringColFilters(colFilters),
		Validator:  validator,
	}

	reader, err := parser.NewParser(file, csvConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating csv: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(reader.Headers())

	for {
		row, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		if err == parser.ErrInvalidRow {
			continue
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unexpected error reading line: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(row)
	}
}

func splitStringColFilters(colFilters string) []string {
	if strings.Trim(colFilters, " ") == "" {
		return nil
	}
	return strings.Split(colFilters, ",")
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", filePath)
		os.Exit(1)
	}
	return file
}

func createRules(rowRulesStr string) []parser.IRule {
	if strings.Trim(rowRulesStr, " ") == "" {
		return nil
	}

	arrRuleStrings := strings.Split(rowRulesStr, ";")

	rules := make([]parser.IRule, 0, len(arrRuleStrings))

	for _, strRule := range arrRuleStrings {
		if strings.Trim(strRule, " ") == "" {
			continue
		}

		r, err := parser.NewRule(strRule)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid rule: %s\n", strRule)
			os.Exit(1)
		}

		rules = append(rules, r)
	}

	return rules
}
