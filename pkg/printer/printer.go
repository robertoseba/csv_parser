package printer

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type Printer struct {
	onScreen    bool
	separator   string
	maxColWidth int
}

func NewPrinter() *Printer {

	return &Printer{
		onScreen:    term.IsTerminal(int(os.Stdout.Fd())),
		separator:   "\t",
		maxColWidth: 30,
	}
}

func (p *Printer) PrintHeader(headers []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(headers, ","))
		return
	}

	fmt.Printf("%s%s\n", p.separator, strings.Join(headers, p.separator))
}

func (p *Printer) PrintRow(lineNumber int, row []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(row, ","))
		return
	}
	fmt.Printf("%d%s%s\n", lineNumber, p.separator, strings.Join(row, p.separator))
}
