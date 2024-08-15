package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/robertoseba/csv_parser/cmd/app"
)

func main() {
	var colFilterFlag = flag.String("filter", "", "Filter the CSV file by the specified columns")
	var colRulesFlag = flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"col1:eq(100)\"")
	flag.Parse()

	filename := flag.Arg(0)

	if filename == "" {
		app.Run(readerStdin(), *colFilterFlag, *colRulesFlag)
		return
	}

	f := readerFile(filename)
	defer f.Close()
	app.Run(f, *colFilterFlag, *colRulesFlag)
}

func readerStdin() *os.File {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Println("No files or pipes provided")
		os.Exit(1)
	}
	return os.Stdin
}

func readerFile(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file: err:", err)
		os.Exit(1)
	}
	return f
}
