package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/robertoseba/csv_parser/cmd/app"
)

func main() {
	var colFilterFlag = flag.String("filter", "", "Filter the CSV file by the specified columns")
	var colRulesFlag = flag.String("rules", "", "Apply rules to the specified columns. Ex: -rules \"col1:eq(100)\"")
	flag.Parse()

	var input io.Reader

	if filename := flag.Arg(0); filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println("error opening file: err:", err)
			os.Exit(1)
		}
		defer f.Close()

		input = f
	} else {
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if fi.Mode()&os.ModeNamedPipe == 0 {
			fmt.Println("No files or pipes provided")
			os.Exit(0)
		}
		input = os.Stdin
	}

	app.Run(input, *colFilterFlag, *colRulesFlag)

}
