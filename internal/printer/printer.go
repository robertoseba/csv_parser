package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Printer struct {
	onScreen    bool
	separator   string
	maxColWidth int
	lineNumber  int
	style       lipgloss.Style
}

func NewPrinter() *Printer {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 3).TabWidth(4)

	return &Printer{
		onScreen:    term.IsTerminal(int(os.Stdout.Fd())),
		separator:   "\t",
		maxColWidth: 0,
		lineNumber:  0,
		style:       baseStyle,
	}
}

func (p *Printer) PrintHeader(headers []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(headers, ","))
		return
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	cellWidth := physicalWidth / len(headers)
	p.style = p.style.MaxWidth(cellWidth)

	p.maxColWidth = cellWidth
	style := p.style.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#D7FF87"))

	for _, header := range headers {
		fmt.Print(style.Render(resizeCell(header, p.maxColWidth)))
	}
	fmt.Println()
	p.lineNumber++
}

func (p *Printer) PrintRow(rows []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(rows, ","))
		return
	}

	style := p.style.
		Foreground(lipgloss.Color("#D7FF87"))

	for _, row := range rows {
		fmt.Print(style.Faint(p.lineNumber%2 == 0).Render(resizeCell(row, p.maxColWidth)))
	}
	fmt.Println()
	p.lineNumber++
}

func resizeCell(cell string, maxWidth int) string {
	if len(cell) > maxWidth {
		return cell[:maxWidth]
	}
	return cell + strings.Repeat(" ", maxWidth-len(cell))
}
