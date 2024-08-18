package printer

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type Printer struct {
	onScreen    bool
	separator   string
	maxColWidth int
	lineNumber  int
	style       lipgloss.Style
	inputChan   <-chan []string
	wg          *sync.WaitGroup
}

func NewPrinter(inputChan chan []string, wg *sync.WaitGroup) *Printer {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().
		Padding(0, 3).
		TabWidth(4)

	return &Printer{
		onScreen:    term.IsTerminal(int(os.Stdout.Fd())),
		separator:   "\t",
		maxColWidth: 0,
		lineNumber:  0,
		style:       baseStyle,
		inputChan:   inputChan,
		wg:          wg,
	}
}

func (p *Printer) Start() {
	startTime := time.Now()
	defer p.terminate(startTime)

	headers := <-p.inputChan
	p.printHeader(headers)

	for rows := range p.inputChan {
		p.printRow(rows)
	}

}

func (p *Printer) printHeader(headers []string) {
	p.calcCellSize(headers)
	headerStyle := p.style.Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#D7FF87"))
	p.print(headerStyle, headers)
}

func (p *Printer) printRow(line []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(line, ","))
		return
	}

	rowStyle := p.style.Foreground(lipgloss.Color("#D7FF87")).Faint(p.lineNumber%2 == 0)
	p.print(rowStyle, line)
}

func (p *Printer) print(style lipgloss.Style, line []string) {
	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(line, ","))
	}

	for _, cell := range line {
		fmt.Print(style.Render(resizeCell(cell, p.maxColWidth)))
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

func (p *Printer) terminate(start time.Time) {
	elapsed := time.Since(start)
	if p.onScreen {
		fmt.Printf("\n Elapsed time: %s\n", elapsed)
		fmt.Println()
	}
	p.wg.Done()
}

func (p *Printer) calcCellSize(headers []string) {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	cellWidth := physicalWidth / len(headers)
	p.style = p.style.MaxWidth(cellWidth)

	p.maxColWidth = cellWidth
}
