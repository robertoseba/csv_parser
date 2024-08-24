package main

import (
	"flag"

	"github.com/robertoseba/csv_parser/cmd/app"
)

func main() {
	options := parseCliOptions()
	app.Run(
		options.filename,
		options.filterInput,
		options.rulesInput,
	)
}

type inputOptions struct {
	filename    string
	filterInput string
	rulesInput  string
}

func parseCliOptions() *inputOptions {
	colFilterFlag := flag.String("filter", "", "Filter the CSV file by the specified columns")
	colRulesFlag := flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"col1:eq(100)\"")
	flag.Parse()

	filename := flag.Arg(0)
	return &inputOptions{
		filename:    filename,
		filterInput: *colFilterFlag,
		rulesInput:  *colRulesFlag,
	}
}

//TODO: Create readme
//TODO: Create CI/CD
//TODO: Publish project
//TODO: Printer calc cell size not based on header but based on first row
//TODO: accept file parameter in different order
