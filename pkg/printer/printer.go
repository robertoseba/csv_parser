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
	lineNumber  int
}

func NewPrinter() *Printer {

	return &Printer{
		onScreen:    term.IsTerminal(int(os.Stdout.Fd())),
		separator:   "\t",
		maxColWidth: 30,
		lineNumber:  0,
	}
}

func (p *Printer) PrintHeader(headers []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(headers, ","))
		return
	}

	fmt.Printf("\x1b[1;30;43m%s%s\x1b[0;m\n", p.separator, strings.Join(headers, p.separator))
	p.lineNumber++
}

func (p *Printer) PrintRow(row []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(row, ","))
		return
	}
	if p.lineNumber%2 == 0 {
		fmt.Printf("\x1b[1;33m%d%s%s\x1b[0;m\n", p.lineNumber, p.separator, strings.Join(row, p.separator))

	} else {
		fmt.Printf("\x1b[2;33m%d%s%s\x1b[0;m\n", p.lineNumber, p.separator, strings.Join(row, p.separator))

	}

	p.lineNumber++
}
