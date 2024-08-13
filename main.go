package main

import (
	"flag"

	"github.com/robertoseba/csv_parser/cmd/app"
)

func main() {
	var colFilterFlag = flag.String("filter", "", "Filter the CSV file by the specified columns")
	var colRulesFlag = flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"col1:eq(100)\"")
	flag.Parse()

	filename := flag.Arg(0)

	// TODO: Implement reading from stdin
	if filename != "" {
		app.Run(filename, *colFilterFlag, *colRulesFlag)
	}

}
