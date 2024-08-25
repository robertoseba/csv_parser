package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/csv_parser/cmd/app"
	"github.com/robertoseba/csv_parser/internal/printer"
)

func main() {
	// TODO: Create readme
	// TODO: Create CI/CD
	// TODO: Publish project
	// TODO: accept file parameter in different order

	inputOptions := app.ParseCliOptions()
	printer := printer.NewPrinter(false)

	err := app.Run(inputOptions, printer)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
