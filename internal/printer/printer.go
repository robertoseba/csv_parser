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
	onScreen         bool
	separator        string
	maxColWidth      []int
	colOverflowAtIdx int
	maxHeight        int
	maxWidth         int
	lineNumber       int
	style            lipgloss.Style
	inputChan        <-chan []string
	wg               *sync.WaitGroup
}

func NewPrinter(inputChan chan []string, wg *sync.WaitGroup) *Printer {
	var width, height int

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 100
		height = 80
	}

	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().
		Padding(0, 2).
		TabWidth(4)

	return &Printer{
		onScreen:   term.IsTerminal(int(os.Stdout.Fd())),
		separator:  "\t",
		maxHeight:  height,
		maxWidth:   width,
		lineNumber: -1,
		style:      baseStyle,
		inputChan:  inputChan,
		wg:         wg,
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
	p.createMaxColWidth(headers)
	// p.calcCellSize(headers)
	headerStyle := p.style.Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#D7FF87"))
	p.print(headerStyle, headers)
}

func (p *Printer) printRow(line []string) {
	rowStyle := p.style.Foreground(lipgloss.Color("#D7FF87")).Faint(p.lineNumber%2 == 0)
	p.print(rowStyle, line)
}

func (p *Printer) print(style lipgloss.Style, line []string) {
	p.lineNumber++

	if !p.onScreen {
		fmt.Printf("%s\n", strings.Join(line, ","))
		return
	}

	if p.lineNumber == 0 { //Printing headers
		fmt.Print(style.Render("Line#"))
	} else {
		fmt.Print(style.Render(fmt.Sprintf("%5d", p.lineNumber)))
	}

	for idx, cell := range line {
		if idx > p.colOverflowAtIdx {
			fmt.Print(style.Foreground(lipgloss.Color("#BBBBBB")).Padding(0).Render("..."))
			break
		}
		fmt.Print(style.Render(resizeCell(cell, p.maxColWidth[idx])))
	}

	fmt.Println()
}

func (p *Printer) terminate(start time.Time) {
	elapsed := time.Since(start)
	if p.onScreen {
		fmt.Println()
		fmt.Printf("%d lines printed\n", p.lineNumber)
		fmt.Printf("Elapsed time: %s\n", elapsed)
		fmt.Println()
	}
	p.wg.Done()
}

func (p *Printer) createMaxColWidth(headers []string) {
	p.maxColWidth = make([]int, 0, len(headers))
	p.colOverflowAtIdx = len(headers)

	hPadding := p.style.GetHorizontalPadding()
	totalWidth := lipgloss.Width("Line#") + hPadding // 6 is the padding

	for idx, header := range headers {
		textWidth := lipgloss.Width(header)
		if totalWidth+textWidth+hPadding > p.maxWidth-lipgloss.Width("...")-hPadding {
			p.colOverflowAtIdx = idx
			textWidth = p.maxWidth - totalWidth - hPadding - lipgloss.Width("...")
			if textWidth < 0 {
				textWidth = 0
			}

			p.maxColWidth = append(p.maxColWidth, textWidth)
			break
		}
		totalWidth += textWidth + hPadding // 6 is the padding
		p.maxColWidth = append(p.maxColWidth, textWidth)
	}
}

func resizeCell(cell string, maxWidth int) string {
	if len(cell) > maxWidth {
		return cell[:maxWidth]
	}
	return cell + strings.Repeat(" ", maxWidth-len(cell))
}
