package app

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/robertoseba/csv_parser/pkg/parser"
)

func Run(filePath string, colFilters string, rowRules string) {

	file, err := os.Open("./data.csv")

	if err != nil {
		fmt.Println("Failed to open file")
		panic(err)
	}

	defer file.Close()

	rules, err := parser.NewRulesFrom(parseRowRules(rowRules))

	if err != nil {
		fmt.Println("Failed to create rules")
		panic(err)
	}

	csvConfig := &parser.CsvConfig{
		Separator:  ',',
		ColFilters: parseColFilters(colFilters),
		Validator:  parser.NewValidator(rules),
	}

	reader, err := parser.New(file, csvConfig)

	if err != nil {
		fmt.Println("Failed to create reader")
		panic(err)
	}

	fmt.Println(reader.FilteredHeaders().Str())

	for {
		row, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Failed to read line")
			panic(err)
		}

		fmt.Println(row.Str())
	}
	return
}

func parseColFilters(colFilters string) []string {
	if strings.Trim(colFilters, " ") == "" {
		return nil
	}
	return strings.Split(colFilters, ",")
}

func parseRowRules(rowRules string) []string {
	if strings.Trim(rowRules, " ") == "" {
		return nil
	}
	return strings.Split(rowRules, "\n")
}
