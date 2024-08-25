package app

import (
	"flag"
	"os"
)

type InputOptions struct {
	Filename    string
	FilterInput string
	RulesInput  string
}

func ParseCliOptions() *InputOptions {
	colFilterFlag := flag.String("filter", "", "Filter the CSV file by the specified columns")
	colRulesFlag := flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"col1:eq(100)\"")

	if os.Args[1][0] == '-' {
		flag.Parse()
	} else {
		// In case the user passes the filename as the first argument before flags, we need to move it to the end
		// so that the flag package can parse the flags correctly
		filename := os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		os.Args = append(os.Args, filename)
		flag.Parse()
	}

	filename := flag.Arg(0)
	return &InputOptions{
		Filename:    filename,
		FilterInput: *colFilterFlag,
		RulesInput:  *colRulesFlag,
	}
}
