package app

import "flag"

type InputOptions struct {
	Filename    string
	FilterInput string
	RulesInput  string
}

func ParseCliOptions() *InputOptions {
	colFilterFlag := flag.String("filter", "", "Filter the CSV file by the specified columns")
	colRulesFlag := flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"col1:eq(100)\"")
	flag.Parse()

	filename := flag.Arg(0)
	return &InputOptions{
		Filename:    filename,
		FilterInput: *colFilterFlag,
		RulesInput:  *colRulesFlag,
	}
}
