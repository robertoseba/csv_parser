package app

import (
	"flag"
	"fmt"
	"os"
)

type InputOptions struct {
	Filename    string
	FilterInput string
	RulesInput  string
	HeaderOnly  bool
}

func ParseCliOptions() *InputOptions {
	colFilterFlag := flag.String("filter", "", "Filter the CSV file by the specified columns from header. Ex: -filter \"username,email\"")
	colRulesFlag := flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"email:eq(user@test.com)\"")
	headerOnlyFlag := flag.Bool("headersOnly", false, "Prints only the headers of the CSV file")

	if len(os.Args) <= 1 {
		wrongUsage()
	}

	sortCliOptions()

	flag.Parse()

	filename := flag.Arg(0)

	if filename == "" {
		wrongUsage()
	}

	return &InputOptions{
		Filename:    filename,
		FilterInput: *colFilterFlag,
		RulesInput:  *colRulesFlag,
		HeaderOnly:  *headerOnlyFlag,
	}
}

func sortCliOptions() {
	if os.Args[1][0] != '-' {
		// In case the user passes the filename as the first argument before flags, we need to move it to the end
		// so that the flag package can parse the flags correctly
		filename := os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		os.Args = append(os.Args, filename)
	}
}

func wrongUsage() {
	fmt.Println("Please specify a csv file to parse and the columns to filter or apply rules.")
	flag.Usage()
	os.Exit(1)
}
