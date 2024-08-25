package main

import (
	"fmt"
	"os"

	"github.com/robertoseba/csv_parser/cmd/app"
	"github.com/robertoseba/csv_parser/internal/printer"
	"golang.org/x/term"
)

func main() {
	inputOptions := app.ParseCliOptions()
	printer := printer.NewPrinter(term.IsTerminal(int(os.Stdout.Fd())))

	err := app.Run(inputOptions, printer)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
